package country

import (
	s "ai30-spatial/pkg/strategy"
	"ai30-spatial/pkg/tools"
	"log"
)

// Choix de l'action :
// Si la planete n'est plus dispo (vide ou autre) alors on abandonne le projet de lancement
// Si une personne à colonisé la map (alors qu'au dernier percept elle ne l'était pas) on utilise les coef de Reaction (soit on se dépêche, soit on ralenti soit on reste pareil)
// Selon la préparation de l'équipage et la stratt, on continue la préparation ou on envoie la fusée
func (p *Country) launchStatus() string {
	highCoef := 4

	// Si la planète est toujorus dispo
	if !tools.ContainsMap(p.available, p.goal) {
		log.Println("Country", p.Name, "is abording the launch to", p.goal.Name, "because the planet is not available anymore")
		p.reset()
		return "abord"
	}

	stratCoef := p.strategyCoef[p.strategy["launch"]]
	distance := p.available[p.goal][0]
	// Si une personne ou plus est sur la planète
	if p.planets[p.goal][2] == 1 {
		stratCoef = tools.GetReaction(stratCoef, p.strategy["reaction"])
	}

	//Si l'équipage est prêt on lance, sinon on continue la préparation
	if p.strategy["launch"] == string(s.EnumLaunch.SLOW) || p.strategy["launch"] == string(s.EnumLaunch.QUICK) || p.strategy["launch"] == string(s.EnumLaunch.MEDIUM) {
		if p.created && (int(distance)/p.rocket.Preparation) <= stratCoef[0] {
			return "launch"
		}
	} else if p.strategy["launch"] == string(s.EnumLaunch.IF_HI_QUICK) || p.strategy["launch"] == string(s.EnumLaunch.IF_HI_SLOW) || p.strategy["launch"] == string(s.EnumLaunch.IF_HI_MEDIUM) {
		cost := p.launchCost(p.goal)
		if cost*highCoef < p.Budget {
			if p.created && (int(distance)/p.rocket.Preparation) < stratCoef[0] {
				return "launch"
			}
		} else if p.created && cost*highCoef > p.Budget && (int(distance)/p.rocket.Preparation) < stratCoef[1] {
			return "launch"
		}
	} else if p.strategy["launch"] == string(s.EnumLaunch.IF_LO_QUICK) || p.strategy["launch"] == string(s.EnumLaunch.IF_LO_SLOW) || p.strategy["launch"] == string(s.EnumLaunch.IF_LO_MEDIUM) {
		cost := p.launchCost(p.goal)
		if cost < p.Budget && cost*2 > p.Budget {
			if p.created && (int(distance)/p.rocket.Preparation) < stratCoef[0] {
				return "launch"
			}
		} else if p.created && cost < p.Budget && (int(distance)/p.rocket.Preparation) < stratCoef[1] {
			return "launch"
		}
	}
	return "prepare"
}
