package rocket

import (
	Planet "ai30-spatial/pkg/planet"
	"ai30-spatial/pkg/tools"
	ws "ai30-spatial/pkg/ws"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

type Rocket struct {
	Id          int
	Name        string
	country     string
	PosX        float64
	PosY        float64
	PosZ        float64
	dest        *Planet.Star
	Velocity    float64
	step        float64
	Preparation int
	chanel      chan string
	ws          *ws.Server
}

func NewRocket(id int, name string, country string, PosX float64, PosY float64, PosZ float64, dest *Planet.Star, Velocity float64, step float64, channel chan string, preparation int, ws *ws.Server) *Rocket {
	return &Rocket{
		Id:          id,
		Name:        name,
		country:     country,
		PosX:        PosX,
		PosY:        PosY,
		PosZ:        PosZ,
		dest:        dest,
		Velocity:    Velocity,
		step:        step,
		chanel:      channel,
		Preparation: preparation,
		ws:          ws,
	}
}

func (r *Rocket) Start() {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	dis := 0.0001
	distance := tools.DistanceBetweenCoordinates(r.PosX, r.PosY, r.PosZ, r.dest.PosX, r.dest.PosY, r.dest.PosZ)
	log.Println("Rocket launched by", r.country, "to", r.dest.Name, "with a preparation of", r.Preparation, " Probability of explosion:", distance/float64(r.Preparation), "distance:", distance)
	dRadius := r.dest.Radius
	tick := r.step
	cnt := 0
	prev := 0
	for (r.PosX < r.dest.PosX-dRadius || r.PosX > r.dest.PosX+dRadius) || (r.PosY < r.dest.PosY || r.PosY > r.dest.PosY+dRadius) || (r.PosZ < r.dest.PosZ || r.PosZ > r.dest.PosZ+dRadius) {
		// log.Println("Rocket is at", r.PosX, r.PosY, r.PosZ)
		// log.Println("Rocket is going to", r.dest.PosX, r.dest.PosY, r.dest.PosZ)
		dirCible := []float64{r.dest.PosX - r.PosX, r.dest.PosY - r.PosY, r.dest.PosZ - r.PosZ}                                   // vecteur directeur
		distance := math.Sqrt(dirCible[0]*dirCible[0] + dirCible[1]*dirCible[1] + dirCible[2]*dirCible[2])                        // Longueur du vecteur (norme)
		direction := []float64{dirCible[0] / distance, dirCible[1] / distance, dirCible[2] / distance}                            // Directon du vecteur
		moving := []float64{direction[0] * r.Velocity * tick, direction[1] * r.Velocity * tick, direction[2] * r.Velocity * tick} // vitesse de déplacement
		newPos := []float64{r.PosX + moving[0], r.PosY + moving[1], r.PosZ + moving[2]}                                           //Nouveau point
		dis += tools.DistanceBetweenCoordinates(r.PosX, r.PosY, r.PosZ, newPos[0], newPos[1], newPos[2])
		r.PosX = newPos[0]
		r.PosY = newPos[1]
		r.PosZ = newPos[2]

		r.ws.Ver.Lock()
		r.ws.WriteMessage([]byte(fmt.Sprintf("{\"id\":%d,\"type\":\"%s\",\"name\":\"%s\", \"sender\":\"%s\",\"x\":%f,\"y\":%f,\"z\":%f,\"status\":\"alive\",\"destination\":\"%s\"}", r.Id, "rocket", "rocket", r.country, r.PosX, r.PosY, r.PosZ, r.dest.Name)))
		r.ws.Ver.Unlock()

		// tick += r.step
		if (int(int(dis) / r.Preparation)) > prev { // Probability of explosion
			cnt++
			prev++
			rand := rand.Intn(100)
			if rand == 50 {
				log.Println("Rocket exploded", "for ", r.country)
				r.ws.Ver.Lock()
				r.ws.WriteMessage([]byte(fmt.Sprintf("{\"id\":%d,\"type\":\"%s\",\"name\":\"%s\", \"sender\":\"%s\",\"x\":%f,\"y\":%f,\"z\":%f, \"status\":\"exploded\",\"destination\":\"%s\"}", r.Id, "rocket", "rocket", r.country, r.PosX, r.PosY, r.PosZ, r.dest.Name)))
				r.ws.Ver.Unlock()
				r.chanel <- "exploded"
				return
			}

		}
		time.Sleep(1 * time.Second / 60)
	}
	r.ws.Ver.Lock()
	r.ws.WriteMessage([]byte(fmt.Sprintf("{\"id\":%d,\"type\":\"%s\",\"name\":\"%s\", \"sender\":\"%s\",\"x\":%f,\"y\":%f,\"z\":%f, \"status\":\"landing\",\"destination\":\"%s\"}", r.Id, "rocket", "rocket", r.country, r.PosX, r.PosY, r.PosZ, r.dest.Name)))
	r.ws.Ver.Unlock()
	log.Println("Rocket arrived to", r.dest.Name, cnt, "random choice mades", "distance:", dis)
	log.Println("Rocket arrived to", r.dest.Name)
	r.chanel <- "arrived"
}
