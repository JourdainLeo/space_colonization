package country

import (
	plt "ai30-spatial/pkg/planet"
	s "ai30-spatial/pkg/strategy"
	"ai30-spatial/pkg/tools"
	"math/rand"
	"time"
)

// pour chaque planète, si il reste des ressources dessus, qu'on a le Budget pour y aller et que selon notre stratégie de skipping on est bon, on l'ajoute aux planètes available (la map qu'on utilise dans les fonctions d'après)
func (p *Country) planetsAvailable() {
	p.available = make(map[*plt.Star][]float64)
	nbEmpty := 0
	countHaveBudget := 0
	//p.writeInLog("je suis ici - %d", p.numberRoundAskAlliance)
	for planet := range p.planets {
		if planet.IsEmpty() {
			nbEmpty++
		}
		if !planet.IsEmpty() && !(planet.IsOwned() && p.strategy["skipping"] == string(s.EnumSkipping.ALWAYS)) && !(planet.IsOwned() && (p.strategy["skipping"] == string(s.EnumSkipping.IF_HI) && p.Budget > 4*p.launchCost(planet))) && !(planet.IsOwned() && (p.strategy["skipping"] == string(s.EnumSkipping.IF_LO) && p.Budget < 2*p.launchCost(planet))) {
			if p.Budget >= p.launchCost(planet) {
				countHaveBudget += 1
				p.available[planet] = make([]float64, 2)
				p.available[planet][0] = tools.DistanceBetweenTwoStar(p.base, planet)
				p.available[planet][1] = float64(p.launchCost(planet))
			} else {
				hasTheMoney := false
				if len(p.alliance) != 0 { // If the country have allies
					part, rest := p.calculateCostForEachAlly(p.launchCost(planet))

					if part != -1 && rest != -1 {
						hasTheMoney = true
						countHaveBudget += 1
						p.available[planet] = make([]float64, 2)
						p.available[planet][0] = tools.DistanceBetweenTwoStar(p.base, planet)
						p.available[planet][1] = float64(p.launchCost(planet))
					}
				}
				if !hasTheMoney {
					flag := true
					loop := 0
					rand.NewSource(time.Now().UnixNano())
					for flag && loop < CMaxLoopToAskAlliance {
						loop += 1
						randomIndex := rand.Intn(len(p.TabCountry))

						if p.askAlliance(p.TabCountry[randomIndex]) {
							flag = false
						}
					}
				}
			}

			p.planets[planet][2] = 0
		} else {
			//Si la map est occupée on note
			if planet.IsOwned() && p.planets[planet][2] == 0 && p.planets[planet][1] != 0 && !tools.Contains(p.ptlAchieved, planet) {
				p.planets[planet][0] = 0
				p.planets[planet][1] = 0
			} else if planet.IsOwned() && p.planets[planet][2] == 0 {
				p.planets[planet][0] /= 2
			}
			p.planets[planet][2] = 1
		}

	}

	if countHaveBudget == 0 && p.digger <= 0 {
		p.numberRoundAskAlliance += 1
	}
	if nbEmpty == len(p.planets) { // Si aucune, fin de la simultaion
		// TODO message de fin de simulation pour le pays
		p.end = true
		return
	} else if len(p.available) == 0 {
		time.Sleep(5 * time.Second)
		p.roundNt++
		if p.roundNt == 10 {
			p.end = true
			return
		}
	}

}
