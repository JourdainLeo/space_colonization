package country

import (
	plt "ai30-spatial/pkg/planet"
	"ai30-spatial/pkg/rocket"
	"ai30-spatial/pkg/tools"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// On envoie la rocket et on attend qu'elle arrive ou crash
// Si elle arrive on commence à miner de la red stone...
// Dans tous les cas on reset les variables de préparation de la fusée et d'objectif de planète
func (p *Country) sendRocket() {

	allyHelpRocket := make([]*Country, 0)
	if !tools.Contains(p.order, p.goal) {
		p.reset()
		p.writeInLog("abording launch")
		return
	}
	p.goal.SetAcquire(p.Id, true)
	if !p.goal.IsSatellite() {
		p.rocket.PosX = p.base.PosX
		p.rocket.PosY = p.base.PosY
		p.rocket.PosZ = p.base.PosZ
	} else {
		p.rocket.PosX = p.base.ParentX
		p.rocket.PosY = p.base.ParentY
		p.rocket.PosZ = p.base.ParentZ
	}
	p.writeToServer("{\"type\":\"country\",\"action\":\"launch\",\"target\":\"" + p.goal.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")

	go p.rocket.Start()

	p.Lock()
	if p.launchCost(p.goal) > p.Budget && p.Budget > CMinMoneyToMakeAlliance {
		somme := 0

		for _, ally := range p.alliance {
			ally.Lock()
		}
		part, rest := p.calculateCostForEachAlly(p.launchCost(p.goal))

		if part == -1 || rest == -1 {
			for _, ally := range p.alliance {
				ally.Unlock()
			}
			p.goal.SetAcquire(p.Id, false)
			p.reset()
			return
		}

		restIsPayed := false
		for _, ally := range p.alliance {
			//fmt.Println("AVANT ", ally.Name, ".budget: ", ally.Budget)
			if ally.Budget > part+rest && !restIsPayed {
				ally.Budget -= part + rest
				allyHelpRocket = append(allyHelpRocket, ally)
				restIsPayed = true
				somme += part + rest
			} else if ally.Budget > part {
				ally.Budget -= part
				allyHelpRocket = append(allyHelpRocket, ally)
				somme += part
			}
			//fmt.Println("APRES ", ally.Name, ".budget: ", ally.Budget)
		}

		for _, ally := range p.alliance {
			ally.Unlock()
		}

		if !restIsPayed { // if rest is not payed

			p.Budget -= rest
		}

		//fmt.Println("AVANT moi ", p.Name, ".budget", p.Budget)
		p.Budget -= (p.launchCost(p.goal) - somme)
		//fmt.Println("APRES moi ", p.Name, ".budget", p.Budget)

	} else {
		p.Budget -= p.launchCost(p.goal)

	}
	p.Unlock()

	for {
		msg := <-p.rocketChnl
		switch msg {
		case "arrived":
			p.nbAchieved++
			p.writeInLog("Rocket arrived to %s", p.goal.Name)
			if p.goal.AddOwner(p.Id, p.Name) {
				p.goal.SetAcquire(p.Id, true)
				var cpy *plt.Star = p.goal
				p.addAchieved(cpy)
				if len(allyHelpRocket) > 0 {
					names := make([]string, 0)
					for _, ally := range allyHelpRocket {
						ally.addAchieved(cpy)
						ally.digger++
						names = append(names, ally.Name)
					}
					p.writeInLog("Digging on %s thanks to [%s]", cpy.Name, strings.Join(names, ", "))
					go p.startDiggerWithAlly(cpy, allyHelpRocket)
				} else {
					p.writeInLog("Digging on %s", cpy.Name)
					go p.startDigger(cpy)
				}

				p.digger++

			} else {
				p.goal.SetAcquire(p.Id, true)
				var cpy *plt.Star = p.goal
				p.addAchieved(cpy)
				p.writeInLog("No resources on %s", p.goal.Name)
			}
			p.reset()
			return
		case "exploded":
			p.writeInLog("Rocket crashed")
			p.writeInLog("%s", p.goal.Name)
			p.goal.SetAcquire(p.Id, false)
			p.reset()
			return
		}
	}
}

// On prépare la fusée, si elle n'a pas encore été créée on la crée
func (p *Country) prepareExploration() {
	if !tools.Contains(p.order, p.goal) {
		p.reset()
		p.writeInLog("abording")
		return
	}
	if !p.created {
		velocity := tools.GetRocketVelocity(p.goal)
		p.rocket = rocket.NewRocket(p.Id, "rocket", p.Name, 0, 0, 0, p.goal, velocity, 0.5, p.rocketChnl, 1, p.goal.Ws)
		p.created = true
	}
	p.writeToServer("{\"type\":\"country\",\"action\":\"preparation\",\"target\":\"" + p.goal.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
	p.rocket.Preparation += 10
	time.Sleep(1 * time.Second)
}

// Observation de la planète selon la stratégie
func (p *Country) observe() {
	//p.writeInLog("je suis la - %d", p.numberRoundAskAlliance)
	switch p.strategy["observation"] {
	case "quick":
		p.observeQuick()
	case "slow":
		p.observeSlow()
	case "if_hi_quick":
		if p.Budget > 50000 {
			p.observeQuick()
		} else {
			p.observeSlow()
		}
	case "if_lo_quick":
		if p.Budget < 20000 {
			p.observeQuick()
		} else {
			p.observeSlow()
		}

	case "if_hi_slow":
		if p.Budget > 50000 {
			p.observeSlow()
		} else {
			p.observeQuick()
		}
	case "if_lo_slow":
		if p.Budget < 20000 {
			p.observeSlow()
		} else {
			p.observeQuick()
		}
	}
}

// Observation rapide : prend la première planète, si elle a déjà été observée 1 fois et qu'on a le Budget pour y aller on y va
func (p *Country) observeQuick() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	p.order = tools.SortPlanet(p.planets, p.available, p.strategy["quick"])
	minVal := tools.MinMapInt(p.planets)
	for _, data := range p.order {
		cost := p.launchCost(data)
		if p.planets[data][1] > 0 {

			if cost < p.Budget {
				p.numberRoundAskAlliance = 0
				p.goal = data
				p.prepare = true
				p.writeToServer("{\"type\":\"country\",\"action\":\"preparation\",\"target\":\"" + data.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
				return
			} else if len(p.alliance) != 0 && p.Budget > CMinMoneyToMakeAlliance {
				part, rest := p.calculateCostForEachAlly(cost)

				if part != -1 && rest != -1 {
					p.numberRoundAskAlliance = 0
					p.goal = data
					p.prepare = true
					p.writeToServer("{\"type\":\"country\",\"action\":\"preparation\",\"target\":\"" + data.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
					return
				}
			}
			if p.Budget < cost {

				rand.NewSource(time.Now().UnixNano())
				flag := true
				loop := 0
				for flag && loop < CMaxLoopToAskAlliance {
					loop += 1
					randomIndex := rand.Intn(len(p.TabCountry))

					if p.askAlliance(p.TabCountry[randomIndex]) { // If a alliance has been sent
						flag = false
					}
				}
			}

		} else if p.planets[data][0] == minVal { //Si pas le Budget alors on observe
			p.planets[data][0] = p.observation(data)
			if tools.Contains(p.ptlAchieved, data) {
				p.planets[data][0] /= 2
			}
			p.planets[data][1]++
			p.writeToServer("{\"type\":\"country\",\"action\":\"observation\",\"target\":\"" + data.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
			return
		}
	}

	if p.digger <= 0 {
		p.numberRoundAskAlliance += 1
	}
	time.Sleep(5 * time.Second)
}

// Observation lente : on observe toutes les planètes 1 fois, puis on choisi la plus rentable. On l'observe une seconde fois, si c'est toujours la plus rentable on y va
func (p *Country) observeSlow() {
	p.order = tools.SortPlanet(p.planets, p.available, p.strategy["slow"])
	for _, data := range p.order {
		if p.planets[data][1] == 0 {
			p.planets[data][0] = p.observation(data)
			p.planets[data][1]++
			p.writeToServer("{\"type\":\"country\",\"action\":\"observation\",\"target\":\"" + data.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
			return
		}
	}
	for _, ind := range p.order {
		if p.planets[ind][1] < 2 { //Si elle a pas été observée au moins 2 fois, on l'observe
			p.planets[ind][0] = p.observation(ind)
			if tools.Contains(p.ptlAchieved, ind) {
				p.planets[ind][0] /= 2
			}
			p.planets[ind][1]++
			p.writeToServer("{\"type\":\"country\",\"action\":\"observation\",\"target\":\"" + ind.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
			return
		} else if p.launchCost(ind) < p.Budget { // Si on a assez de Budget, on y va
			p.goal = ind
			p.prepare = true
			p.numberRoundAskAlliance = 0
			p.writeToServer("{\"type\":\"country\",\"action\":\"preparation\",\"target\":\"" + ind.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
			return
		} else if len(p.alliance) != 0 && p.Budget > CMinMoneyToMakeAlliance { // If the country have allies
			part, rest := p.calculateCostForEachAlly(p.launchCost(ind))

			if part != -1 && rest != -1 { // if the country with his ally have the money
				p.goal = ind
				p.prepare = true
				p.numberRoundAskAlliance = 0
				p.writeToServer("{\"type\":\"country\",\"action\":\"preparation\",\"target\":\"" + ind.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
				return
			}
		}
		if p.Budget < p.launchCost(ind) {

			rand.NewSource(time.Now().UnixNano())
			flag := true
			loop := 0
			for flag && loop < CMaxLoopToAskAlliance {
				loop += 1
				randomIndex := rand.Intn(len(p.TabCountry))

				if p.askAlliance(p.TabCountry[randomIndex]) { // If a alliance has been sent
					flag = false
				}
			}

		}
	}
	if p.digger <= 0 {
		p.numberRoundAskAlliance += 1
	}
	time.Sleep(5 * time.Second)
}

func (p *Country) calculateCostForEachAlly(cost int) (int, int) {

	if len(p.alliance) <= 0 {
		return -1, -1
	}

	diff := cost - p.Budget + CMinMoneyToMakeAlliance

	rest := diff % len(p.alliance)
	part := diff / len(p.alliance)

	//fmt.Println("diff: ", diff, "rest:", rest, "part:", part)

	cantPay := 0
	canPayRest := false
	isNOK := true

	maxLoop := len(p.alliance) + 10
	loop := 0

	for isNOK && cantPay < len(p.alliance) && loop < maxLoop {
		loop++
		//p.writeInLog("Im here")
		rest = diff % (len(p.alliance) - cantPay)
		part = diff / (len(p.alliance) - cantPay)
		isNOK = false
		for _, ally := range p.alliance {

			if ally.Budget > part+rest {
				canPayRest = true

			} else if ally.Budget <= part {
				cantPay -= 1

				if cantPay < 0 {
					isNOK = true
				}
			}
		}

		if !canPayRest {
			cantPay -= 1
			isNOK = true
		}

		if isNOK {
			cantPay *= -1
			canPayRest = false
		}
	}

	if loop >= maxLoop {
		return -1, -1
	}

	if isNOK {
		return -1, -1
	}

	return part, rest

}
