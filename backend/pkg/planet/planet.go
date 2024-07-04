package Planet

import (
	ws "ai30-spatial/pkg/ws"
	"fmt"
	"log"
	"sync"
)

type Planet struct {
	Star
	satellites []*Satellite
}

func NewPlanet(id int, name string, resources int, velocity float64, posX float64, posY float64, posZ float64, counter *float64, min int, max int, parentSize float64, radius float64, ws *ws.Server, wg *sync.WaitGroup, str *string, srv *sync.Mutex) *Planet {
	return &Planet{
		Star:       *NewStar(id, name, resources, velocity, posX, posY, posZ, counter, min, max, parentSize, radius, *ws, wg, str, srv),
		satellites: make([]*Satellite, 0),
	}
}

func (p *Planet) AddSatellite(sat *Satellite) {
	if !p.Contains(p.Id) {
		p.Lock()
		defer p.Unlock()
		p.satellites = append(p.satellites, sat)

		log.Printf("Le satellite \"%s\" (n°%d) a été ajouté à la planète \"%s\" (n°%d)\n", sat.Name, sat.Id, p.Name, p.Id)
	}
}

func (p *Planet) Contains(idSatellite int) bool {
	return p.GetSatelliteById(idSatellite) != nil
}

func (p *Planet) GetSatelliteById(idSatellite int) *Satellite {
	p.Lock()
	defer p.Unlock()
	for _, sat := range p.satellites {
		if sat.Id == idSatellite {
			return sat
		}
	}
	return nil
}
func (star *Planet) Start() {
	star.Move()
	// log.Println("Position of", star.Name, ":", star.PosX, star.PosY, star.PosZ)
	star.srv.Lock()
	*star.str += fmt.Sprintf("{\"id\":%d,\"type\":\"%s\",\"name\":\"%s\",\"x\":%f,\"y\":%f,\"z\":%f, \"resources\":\"%d\", %s}, ", star.Id, "planet", star.Name, star.PosX, star.PosY, star.PosZ, star.actual, star.getOwnersNames())
	star.srv.Unlock()
	star.wg.Done()
}
