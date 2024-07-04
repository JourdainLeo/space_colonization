package solar_system

import (
	plt "ai30-spatial/pkg/planet"
	"log"
	"sync"
)

type SolarSystem struct {
	sync.Mutex
	planets        []*plt.Planet
	listEmptyIndex []int
	ids            int
}

func NewSolarSystem() *SolarSystem {
	return &SolarSystem{sync.Mutex{}, make([]*plt.Planet, 0), make([]int, 0), 0}
}

func (system *SolarSystem) GetNumberPlanet() int {
	system.Lock()
	defer system.Unlock()
	return len(system.planets)
}

func (system *SolarSystem) GetPlanetById(idPlanet int) *plt.Planet {
	system.Lock()
	defer system.Unlock()
	for _, planet := range system.planets {

		if planet.Id == idPlanet {
			return planet
		}
	}
	return nil
}

func (system *SolarSystem) GetPlanets() []*plt.Planet {
	system.Lock()
	defer system.Unlock()
	return system.planets
}

func (system *SolarSystem) Contains(idPlanet int) bool {
	return system.GetPlanetById(idPlanet) != nil
}

func (system *SolarSystem) AddPlanet(p *plt.Planet) {
	if !system.Contains(p.Id) {
		system.Lock()
		defer system.Unlock()
		if len(system.listEmptyIndex) != 0 { // si un index n'est pas utilisé
			index := 0
			index, system.listEmptyIndex = system.listEmptyIndex[0], system.listEmptyIndex[1:]
			system.planets[index] = p
		} else { // sinon, on ajoute de la place dans le slice
			system.planets = append(system.planets, p)
		}

		log.Printf("La planet \"%s\" (n°%d) a été ajouté à la liste\n", p.Name, p.Id)
	}
}

func (system *SolarSystem) RemovePlanet(idPlanet int) {
	system.Lock()
	defer system.Unlock()
	for index, planet := range system.planets {

		if planet.Id == idPlanet { // si l'id de la planète est trouvé
			system.planets[index] = nil
			system.listEmptyIndex = append(system.listEmptyIndex, index)
			log.Printf("La planet \"%s\" (n°%d) a été enlevé de la liste\n", planet.Name, planet.Id)
			return
		}
	}

}
func (sol *SolarSystem) GetNewId() int {
	sol.Lock()
	defer sol.Unlock()
	sol.ids += 1
	return sol.ids
}
