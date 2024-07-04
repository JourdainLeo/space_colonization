package Planet

import (
	ws "ai30-spatial/pkg/ws"
	"fmt"
	"sync"
)

type Satellite struct {
	Star
	Parent *Planet
}

func NewSatellite(id int, name string, resources int, velocity float64, posX float64, posY float64, posZ float64, counter *float64, min int, max int, parentSize float64, radius float64, ws *ws.Server, parent *Planet, wg *sync.WaitGroup, str *string, srv *sync.Mutex) *Satellite {

	return &Satellite{
		Star:   *NewStar(id, name, resources, velocity, posX, posY, posZ, counter, min, max, parentSize, radius, *ws, wg, str, srv),
		Parent: parent,
	}
}
func (star *Satellite) StartSat() {
	star.Move()
	star.ParentX = star.Parent.PosX
	star.ParentY = star.Parent.PosY
	star.ParentZ = star.Parent.PosZ
	star.srv.Lock()
	*star.str += fmt.Sprintf("{\"id\":%d,\"type\":\"%s\",\"name\":\"%s\",\"x\":%f,\"y\":%f,\"z\":%f, \"resources\":\"%d\", %s}, ", star.Id, "satellite", star.Name, star.PosX, star.PosY, star.PosZ, star.actual, star.getOwnersNames())
	star.srv.Unlock()
	star.wg.Done()
}
