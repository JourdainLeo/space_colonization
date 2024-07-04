package country

import (
	plt "ai30-spatial/pkg/planet"
	"ai30-spatial/pkg/rocket"
	s "ai30-spatial/pkg/strategy"
	"ai30-spatial/pkg/tools"
	ws "ai30-spatial/pkg/ws"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

const CWaitForAlliance = 5
const CMinMoneyToMakeAlliance = 1000
const CLimitAskAlliance = 10
const CMaxLoopToAskAlliance = 20

type AskedAlliance struct {
	IdSender         int
	IdReceiver       int
	InfoOtherCountry *Country
	Time             int
}

type Country struct {
	sync.Mutex
	watch        sync.Mutex
	Id           int
	Name         string
	nbScientists []int
	iq           int
	Budget       int
	nbAchieved   int
	nbParts      float32
	alliance     []*Country
	// list of pair {IdSender, IdReceive, NumberDays} the country asked
	// NumberDays : count how many days left before the alliance is expired
	allianceInWait []AskedAlliance
	//Map de panète dont l'ordre ne change pas, on actualise juste les ressources de cette planète et le nombre de fois qu'on l'a observé
	planets      map[*plt.Star][]int // Map of planets, first index of tab is the ressources observed and the second is how manu time it has been observed, third is if the planet is owned
	goal         *plt.Star
	base         *plt.Star
	rocket       *rocket.Rocket
	channels     []chan interface{}
	strategy     map[string]string
	strategyCoef map[string][]int
	pltStrategy  map[*plt.Star]string
	//Map de planète dont l'ordre change selon la stratégie, on actualise la distance et le prix de lancement
	available              map[*plt.Star][]float64 // Map of planets -> distance, price
	order                  []*plt.Star             //Order of observation of the planets
	end                    bool
	created                bool
	prepare                bool
	rocketChnl             chan string
	Ws                     *ws.Server
	TabCountry             []*Country
	ptlAchieved            []*plt.Star
	digger                 int
	wg                     (*sync.WaitGroup)
	numberRoundAskAlliance int
	roundNt                int
}

func NewCountry(id int, name string, nbScientists []int, Budget int, nbAchieved int, alliance []*Country, planets map[*plt.Star][]int, goal *plt.Star, base *plt.Star, rocket *rocket.Rocket, channels []chan interface{}, strategy map[string]string, strategyCoef map[string][]int, order []*plt.Star, ws *ws.Server, wg *sync.WaitGroup) *Country {
	return &Country{
		watch:                  sync.Mutex{},
		Id:                     id,
		Name:                   name,
		nbScientists:           nbScientists,
		iq:                     getIq(nbScientists),
		Budget:                 Budget,
		nbAchieved:             nbAchieved,
		nbParts:                0,
		alliance:               alliance,
		planets:                planets,
		goal:                   goal,
		base:                   base,
		rocket:                 rocket,
		channels:               channels,
		strategy:               strategy,
		strategyCoef:           strategyCoef,
		pltStrategy:            make(map[*plt.Star]string),
		end:                    false,
		created:                false,
		prepare:                false,
		rocketChnl:             make(chan string),
		available:              make(map[*plt.Star][]float64),
		order:                  order,
		Ws:                     ws,
		ptlAchieved:            make([]*plt.Star, 0),
		digger:                 0,
		wg:                     wg,
		roundNt:                0,
		numberRoundAskAlliance: 0,
	}
}

// Return the sum of all the iq of the scientist from the country
func getIq(nbScientists []int) int {
	sum := 0
	for _, iq := range nbScientists {
		sum += iq
	}
	return sum
}

// observation penalty
func (p *Country) observation(plt *plt.Star) int {
	dist := tools.DistanceBetweenTwoStar(p.base, plt)
	p.writeToServer("{\"type\":\"country\",\"action\":\"observation\",\"target\":\"" + plt.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
	dur := (int(dist) / 4000) + 1
	time.Sleep(time.Duration(dur) * time.Second)
	return plt.ObserveResources(p.iq, p.Id)
}

// Launch price : 1000
func (p *Country) launchCost(planet *plt.Star) int {
	dis := tools.DistanceBetweenTwoStar(p.base, planet)
	return int(dis + 10000)
}

// Dig planet
func (p *Country) startDigger(planet *plt.Star) {
	for {
		log.Println("Country", p.Name, "is digging on", planet.Name)
		val := planet.Dig(p.Id)
		if val > 0 {
			p.Lock()
			p.Budget += val
			p.nbParts += 10
			p.Unlock()
		} else {
			p.digger--
			return
		}
		time.Sleep(4 * time.Second)
	}
}

// Dig planet with help
func (p *Country) startDiggerWithAlly(planet *plt.Star, allies []*Country) {
	for {
		val := planet.Dig(p.Id)
		if val > 0 {

			partPlanetForEach := 10 / (len(allies) + 1)
			restPlanetForEach := 10 % (len(allies) + 1)
			part := val / (len(allies) + 1)
			rest := val % (len(allies) + 1)

			for _, ally := range allies {
				ally.Lock()
				ally.Budget += part
				ally.nbParts += float32(partPlanetForEach)
				ally.Unlock()
			}
			p.Lock()
			p.Budget += part + rest
			p.nbParts += float32(partPlanetForEach) + float32(restPlanetForEach)
			p.Unlock()
		} else {
			p.digger--
			for _, ally := range allies {
				ally.digger--
			}
			return
		}
		log.Println("Country", p.Name, "is digging on", planet.Name, "with allies")
		time.Sleep(4 * time.Second)

	}
}

// Accept or refuse an alliance
func (p *Country) processAlliance(pAlliance *Country) (response bool) {

	if p.strategy["alliance"] == string(s.EnumAlliance.NEVER) {
		response = false
	} else if p.strategy["alliance"] == string(s.EnumAlliance.ALWAYS) {

		response = true
	} else if p.strategy["alliance"] == string(s.EnumAlliance.IF_HI) {

		if p.Budget >= 50000 {
			response = true
		} else {
			response = false
		}
	} else if p.strategy["alliance"] == string(s.EnumAlliance.IF_LO) {

		if p.Budget <= 20000 {
			response = true
		} else {
			response = false
		}

	}

	if response {
		p.alliance = append(p.alliance, pAlliance)
		p.writeInLog("Reception of request of country %s (n°%d) for a alliance and has been accepted", pAlliance.Name, pAlliance.Id)
		p.writeToServer("{\"type\":\"country\",\"action\":\"alliance\",\"target\":\"" + pAlliance.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
	} else {
		p.writeInLog("Reception of request of country %s (n°%d) for a alliance and has been REJECTED", pAlliance.Name, pAlliance.Id)
	}

	return
}

// treat all messages of alliance received until one is accepted
// retur `True` if a alliance is accepted
func (p *Country) treatAlliance() bool {
	flag := false
	index := 0
	p.Lock()
	for len(p.allianceInWait) > index && !flag {

		if p.allianceInWait[index].IdSender == p.Id { // If this country has proposed a alliance

			p.allianceInWait[index].Time--

			if p.allianceInWait[index].Time <= 0 { // If we finished to wait
				p.writeInLog("Cancel alliance with the country %s (n°%d) ", p.allianceInWait[index].InfoOtherCountry.Name, p.allianceInWait[index].InfoOtherCountry.Id)

				p.allianceInWait = slices.Delete(p.allianceInWait, index, index+1) // remove the proposition

			} else {
				index++
			}

		} else { // If a other country has proposed a alliance

			countryToRespond := p.allianceInWait[index].InfoOtherCountry

			if !slices.ContainsFunc(p.alliance, func(e *Country) bool { // if the country is not a friend yet

				return e.Id == countryToRespond.Id

			}) {

				flag = p.processAlliance(countryToRespond) // do we accept or not
				if flag {
					idToRespond := countryToRespond.Id
					p.channels[idToRespond] <- MsgAllianceRespond{p}
				}
			}

			p.allianceInWait = slices.Delete(p.allianceInWait, index, index+1) // remove the proposition

		}

	}
	p.Unlock()

	return flag
}

// Listen messages from another country
func (p *Country) listenAlliance() {

	flag := true

	for !p.end && flag {

		select {

		case msg := <-p.channels[p.Id]:

			switch m := msg.(type) {

			case MsgAlliance:

				p.Lock()
				p.allianceInWait = append(p.allianceInWait, AskedAlliance{m.Sender.Id, p.Id, m.Sender, CWaitForAlliance})
				p.Unlock()

			case MsgAllianceRespond:
				p.Lock()

				index := slices.IndexFunc(p.allianceInWait, func(e AskedAlliance) bool {
					return e.IdReceiver == m.Sender.Id
				})

				if index >= 0 {

					p.allianceInWait = append(p.allianceInWait[:index], p.allianceInWait[index+1:]...)
					p.alliance = append(p.alliance, m.Sender)
					p.writeToServer("{\"type\":\"country\",\"action\":\"alliance\",\"target\":\"" + m.Sender.Name + "\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
					p.writeInLog("Send request to country %s (n°%d) for a alliance and has been accepted", m.Sender.Name, m.Sender.Id)
				}
				p.Unlock()

				//p.wgAlliance.Done()

			default:

			}

		default:
			flag = false

		}
	}
}

func (p *Country) writeInLog(msg string, parameter ...interface{}) {
	textToAdd := fmt.Sprintf(msg, parameter...)
	log.Printf("Country \"%s\" (n°%d) - %s\n", p.Name, p.Id, textToAdd)
}

func (p *Country) askAlliance(friend *Country) bool {

	p.Lock()
	defer p.Unlock()
	if slices.ContainsFunc(p.alliance, func(e *Country) bool { // if the country is already a friend

		return e.Id == friend.Id

	}) {
		return false
	}

	if slices.ContainsFunc(p.allianceInWait, func(e AskedAlliance) bool {

		return e.IdReceiver == friend.Id || e.IdSender == friend.Id
	}) {
		return false
	}

	if p.Id == friend.Id { // if the country has the same id
		return false
	}

	p.allianceInWait = append(p.allianceInWait, AskedAlliance{p.Id, friend.Id, friend, CWaitForAlliance})

	select {
	case p.channels[friend.Id] <- MsgAlliance{p}:
		p.writeInLog("Ask a alliance to %s (n°%d) ", friend.Name, friend.Id)

	default:
		return false
	}

	return true
}

func (p *Country) reset() {
	p.prepare = false
	p.created = false
	p.goal = &plt.Star{}
}

func (p *Country) addAchieved(plt *plt.Star) {
	if !slices.Contains(p.ptlAchieved, plt) {
		p.ptlAchieved = append(p.ptlAchieved, plt)
		p.nbAchieved++
	}
}

// Observation de l'environnement
func (p *Country) percept() {
	p.listenAlliance()
	p.planetsAvailable()
	// Sort planets depend on strat
	// p.listenAlliance() // à voir la compatibilité
	// Observe uniquement l'environnement c'est à dire on enregistre les données courantes -> fusée en cours, distance aux planètes / satellites, ressources en cours, les nvl demandes d'alliances (imo) etc...
}

// Choix de l'action : Observation, préparation de la fusée, lancement, abandon de la préparation
func (p *Country) deliberate() string {
	// selon notre stratégie, l'état courant donné par percept, nos ressources, notre QI, nos alliance ou autre -> on renvoie un string / enum selon les meilleures actions
	if p.prepare {
		return p.launchStatus()
	}
	return "observe"
	// Ajouter condition pour ask alliance
}

// Réalisation de l'action choisie
func (p *Country) act(str string) {
	// on réalise l'action : selon le string / enum renvoyé par deliberate(), on appelle la fonction qui fait l'action voulu

	if p.treatAlliance() { // if a new alliance is accepted
		return
	}
	if str == "observe" {
		//	log.Println("Country", p.Name, "is observing")
		p.observe()
	} else if str == "prepare" {
		//	log.Println("Country", p.Name, "is preparing")
		p.prepareExploration()
	} else if str == "launch" {
		//	log.Println("Country", p.Name, "is launching")
		p.sendRocket()
	} else if str == "abord" {
		//	log.Println("Country", p.Name, "is abording")
		p.reset()
	}

}

func (p *Country) getAlliesNames() string {
	names := make([]string, 0)
	for _, ally := range p.alliance {
		names = append(names, ally.Name)
	}
	return strings.Join(names, ", ")
}

func (p *Country) GetSummary() string {
	str := fmt.Sprintf("{\"country\":\"%s\", \"presence\":\"%d\", \"planets\":\"%.1f\", \"budget\":\"%d\", \"allies\":\"%s\", \"strategy\":[{\"observation\":\"%s\",\"launch\":\"%s\",\"skipping\":\"%s\", \"alliance\":\"%s\"}]}", p.Name, p.nbAchieved, float32(p.nbParts/100), p.Budget, p.getAlliesNames(), p.strategy["observation"], p.strategy["launch"], p.strategy["skipping"], p.strategy["alliance"])
	log.Println(str)
	return str
}

// Start the loop for the agent
func (p *Country) Start() {
	//go p.listenAlliance()
	// If no goal or rocket sent, observe.
	p.writeToServer("{\"type\":\"country\",\"action\":\"perception\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + ",\"qi\":" + strconv.Itoa(p.iq) + "}")
	for !p.end && p.numberRoundAskAlliance < CLimitAskAlliance {
		p.percept()
		p.act(p.deliberate())
		// p.writeToServer()
	}
	p.end = true
	p.writeToServer("{\"type\":\"country\",\"action\":\"ended\",\"country\":\"" + p.Name + "\", \"budget\":" + strconv.Itoa(p.Budget) + ", " + tools.GetPlanetName(p.ptlAchieved) + "}")
	p.writeInLog("Simulation ended")
	p.writeInLog("Has achieved %d planets and has %d Budget", p.nbAchieved, p.Budget)
	p.wg.Done()
	// else prepare launch
}

func (p *Country) writeToServer(message string) {
	p.Ws.Ver.Lock()
	p.Ws.WriteMessage([]byte(message))
	p.Ws.Ver.Unlock()
}
