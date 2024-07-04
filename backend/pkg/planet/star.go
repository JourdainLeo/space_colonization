package Planet

import (
	ws "ai30-spatial/pkg/ws"
	"log"
	"math"
	"math/rand"
	"slices"
	"sync"
	"time"
)

type Star struct {
	sync.Mutex
	Id          int
	Name        string
	acquired    map[int]bool
	PosX        float64
	PosY        float64
	PosZ        float64
	ParentX     float64
	ParentY     float64
	ParentZ     float64
	observed    map[int]int // Map of countries who observed the planet, the int is the number of time it has been observed
	owner       int
	owners      map[int]int // Map of countries who own the planet, the int is the percentage of the planet owned
	ownersName  []string
	resources   int
	actual      int
	velocity    float64
	maxDistance int
	minDistance int
	counter     *float64
	parentSize  float64
	Radius      float64
	wg          *sync.WaitGroup
	Ws          *ws.Server
	str         *string
	srv         *sync.Mutex
}

func NewStar(id int, name string, resources int, velocity float64, posX float64, posY float64, posZ float64, counter *float64, min int, max int, parentSize float64, radius float64, ws ws.Server, wg *sync.WaitGroup, str *string, srv *sync.Mutex) *Star {
	return &Star{
		Id:          id,
		Name:        name,
		acquired:    make(map[int]bool),
		resources:   resources,
		actual:      resources,
		velocity:    velocity,
		PosX:        posX,
		PosY:        posY,
		PosZ:        posZ,
		ParentX:     -1,
		ParentY:     -1,
		ParentZ:     -1,
		owner:       0,
		owners:      make(map[int]int),
		ownersName:  make([]string, 0),
		maxDistance: max,
		minDistance: min,
		parentSize:  parentSize,
		counter:     counter,
		Radius:      radius,
		observed:    make(map[int]int),
		Ws:          &ws,
		wg:          wg,
		str:         str,
		srv:         srv,
	}
}

func (star *Star) Colonised(idPays int) int {
	star.Lock()
	defer star.Unlock()

	if star.IsAcquired() {
		return -1
	} else {
		star.acquired[idPays] = true
		star.owner = idPays
		log.Printf("Star \"%s\" (n°%d) has been colonised by (n°%d)\n", star.Name, star.Id, idPays)
		return idPays
	}
}

// SetAcquire : Set the planet as acquired by the country
func (star *Star) IsAcquired() bool {
	star.Lock()
	defer star.Unlock()
	for _, v := range star.acquired {
		if v {
			return true
		}
	}
	return false
}

func (star *Star) SetAcquire(idPays int, state bool) {
	star.Lock()
	defer star.Unlock()
	star.acquired[idPays] = state
}

func (star *Star) IsOwned() bool {
	star.Lock()
	defer star.Unlock()
	return len(star.owners) > 0
}

func (star *Star) AddOwner(idPays int, name string) bool {
	star.Lock()
	defer star.Unlock()
	if star.owner < 100 {
		star.owners[idPays] = 0
		if !slices.Contains(star.ownersName, name) {
			star.ownersName = append(star.ownersName, name)
		}
		return true
	}
	return false
}

func (star *Star) Dig(id int) int {
	star.Lock()
	defer star.Unlock()
	if star.owner < 100 {
		star.owner += 10
		star.actual -= star.resources / 10
		return star.resources / 10
	}
	return 0
}

func (star *Star) IsEmpty() bool {
	star.Lock()
	defer star.Unlock()
	return star.owner == 100
}

// ObserveResources : Depends on the IQ and the number of time that the country observed -> observation is more precise or not
func (star *Star) ObserveResources(qi int, pays int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	star.Lock()
	defer star.Unlock()
	if _, ok := star.observed[pays]; !ok {
		star.observed[pays] = 0
	}
	base := int(5000 / (star.observed[pays] + 1))
	minRes := star.actual - base + qi
	maxRes := star.actual + base + qi
	ret := rand.Intn(maxRes-minRes+1) + minRes

	if star.observed[pays] == 0 {
		star.observed[pays] = 1
	} else {
		star.observed[pays]++
	}
	return ret
}

func (star *Star) GetPosition() (float64, float64, float64) {
	star.Lock()
	defer star.Unlock()
	return star.PosX, star.PosY, star.PosZ
}

func (star *Star) getOwnersNames() string {
	star.Lock()
	defer star.Unlock()
	explo := "\"owners\": ["
	added := false
	for _, v := range star.ownersName {
		explo += "\"" + v + "\", "
		added = true
	}
	if added {
		explo = explo[:len(explo)-2]
	}
	explo += "]"
	return explo

}

func (star *Star) IsSatellite() bool {
	return star.ParentX != -1 && star.ParentY != -1 && star.ParentZ != -1
}

func (star *Star) Move() {
	star.Lock()
	defer star.Unlock()
	cnt := *star.counter
	var angle = float64(cnt * star.velocity)
	star.PosX = math.Sin(angle) * (float64(star.maxDistance) + star.parentSize)
	star.PosZ = math.Cos(angle) * (float64(star.minDistance) + star.parentSize)

	// fmt.Printf("Star n°%d %s: %f, %f, %f\n", star.Id, star.Name, star.PosX, star.PosY, star.PosZ)
}
