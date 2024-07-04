package main

import (
	"ai30-spatial/pkg/country"
	plt "ai30-spatial/pkg/planet"
	"ai30-spatial/pkg/tools"
	ws "ai30-spatial/pkg/ws"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type PlanetData struct {
	Velocity    float64  `json:"velocity"`
	Size        float64  `json:"size"`
	MinDistance int      `json:"minDistance"`
	MaxDistance int      `json:"maxDistance"`
	Satellites  []string `json:"satellites"`
}

type SatelliteData struct {
	Velocity float64 `json:"velocity"`
	Distance int     `json:"distance"`
}

type SystemData struct {
	Planets    map[string]PlanetData    `json:"planets"`
	Satellites map[string]SatelliteData `json:"satellites"`
}
type start struct {
	Pays map[string]int `json:"countries"`
}

func main() {
	server := ws.StartServer()
	var cnt float64 = 0.016
	wg := new(sync.WaitGroup)
	var str string
	var srv = new(sync.Mutex)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Open JSON
	file, err := os.Open("./data.json")
	if err != nil {
		fmt.Println("Error opening the json :", err)
		return
	}
	defer file.Close()

	// Decode JSON
	var system SystemData
	err = json.NewDecoder(file).Decode(&system)
	if err != nil {
		fmt.Println("Error decoding the json:", err)
		return
	}

	var p []*plt.Planet
	var s []*plt.Satellite
	var planets map[*plt.Star][]int
	var countries []*country.Country
	var base *plt.Star
	planets = make(map[*plt.Star][]int)
	strategies := make(map[string][]string)
	strategiesCoef := make(map[string][]int)
	strategies["observation"] = []string{"quick", "slow", "if_hi_quick", "if_hi_slow", "if_lo_quick", "if_lo_slow"}
	// strategies["observation"] = []string{"slow"}
	strategies["launch"] = []string{"quick", "slow", "medium", "if_hi_quick", "if_hi_slow", "if_hi_medium", "if_lo_quick", "if_lo_slow", "if_lo_medium"}
	strategies["alliance"] = []string{"never", "always", "if_hi", "if_lo"}
	strategies["skipping"] = []string{"never", "always", "if_hi", "if_lo"}
	strategies["reaction"] = []string{"faster", "slower", "same", "balanced_hi", "balanced_lo"}
	strategiesCoef["quick"] = []int{50, 50}
	strategiesCoef["slow"] = []int{5, 5}
	strategiesCoef["medium"] = []int{25, 25}
	strategiesCoef["if_hi_quick"] = []int{50, 5}
	strategiesCoef["if_lo_quick"] = []int{50, 5}
	strategiesCoef["if_lo_medium"] = []int{25, 50}
	strategiesCoef["if_hi_medium"] = []int{25, 50}
	strategiesCoef["if_hi_slow"] = []int{5, 50}
	strategiesCoef["if_lo_slow"] = []int{5, 50}

	id := 0
	posX, posY, posZ := 0.0, 0.0, 0.0
	for name, value := range system.Planets {
		planet := plt.NewPlanet(id, name, rand.Intn(40000)+10000, value.Velocity, posX, posY, posZ, &cnt, value.MinDistance, value.MaxDistance, 100, 10, server, wg, &str, srv)
		p = append(p, planet)
		if planet.Name == "earth" {
			base = &planet.Star
		} else {
			planets[&planet.Star] = []int{0, 0, 0}
		}

		// Start satellite
		if len(value.Satellites) > 0 {
			for _, sat := range value.Satellites {
				satellite := plt.NewSatellite(id, sat, rand.Intn(20000)+10000, system.Satellites[sat].Velocity, posX, posY, posZ, &cnt, system.Satellites[sat].Distance, system.Satellites[sat].Distance, value.Size, 50, server, planet, wg, &str, srv)
				planet.AddSatellite(satellite)
				s = append(s, satellite)
				planets[&satellite.Star] = []int{0, 0, 0}
			}
		}
		id++
	}

	// Create countries
	id = 0
	fmt.Println("Waiting for countries")
	params := <-server.Channel
	//params := []byte("{\"countries\": {\"france\": 27073,\"north_korea\": 95266,\"south_korea\": 96458,\"united_states\": 95941,\"united_kingdom\": 312,\"india\": 37842,\"china\": 5686,\"russia\": 69374,\"central_africa\": 3999,\"australia\": 92548,\"italy\": 89710,\"germany\": 66063,\"mexico\": 66511,\"canada\": 25904,\"japan\": 77175}}")
	var mapCountry start
	_ = json.Unmarshal(params, &mapCountry)
	numberCountry := len(mapCountry.Pays)
	countrywg := new(sync.WaitGroup)
	channels := make([]chan interface{}, numberCountry)
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan interface{}, numberCountry*2)
	}
	for cntry, ress := range mapCountry.Pays {
		nbScientists := rand.Intn(5)
		scientists := []int{}
		for j := 0; j < nbScientists; j++ {
			scientists = append(scientists, rand.Intn(80)+100)
		}
		strat := make(map[string]string)
		strat["observation"] = strategies["observation"][rand.Intn(len(strategies["observation"]))]
		strat["launch"] = strategies["launch"][rand.Intn(len(strategies["launch"]))]
		strat["alliance"] = strategies["alliance"][rand.Intn(len(strategies["alliance"]))]
		strat["skipping"] = strategies["skipping"][rand.Intn(len(strategies["skipping"]))]
		strat["reaction"] = strategies["reaction"][rand.Intn(len(strategies["reaction"]))]
		order := tools.ShufflePlanet(planets)
		countrywg.Add(1)
		country := country.NewCountry(id, cntry, scientists, ress, 0, []*country.Country{}, planets, nil, base, nil, channels, strat, strategiesCoef, order, server, countrywg)
		countries = append(countries, country)
		id++
	}

	end := make(chan bool)
	go func() {
		for _, country := range countries {
			country.TabCountry = countries
			log.Println("Start country", country.Id)
			go country.Start()
		}
		countrywg.Wait()
		end <- true
	}()
	finish := false
	// Start planet
	for !finish {
		cnt += 0.016
		srv.Lock()
		str = "{\"type\":\"star\", \"data\":["
		srv.Unlock()
		wg.Add(len(p) + len(s))
		for _, pl := range p {
			go pl.Start()
		}
		// Start satellite
		for _, sat := range s {
			go sat.StartSat()
		}
		wg.Wait()
		str = str[:len(str)-2]
		str += "]}"
		server.Ver.Lock()
		server.WriteMessage([]byte(str))
		server.Ver.Unlock()
		time.Sleep(1 * time.Second / 60)
		select {
		case <-end:
			finish = true
		default:
			finish = false

		}
	}
	sum := "{\"type\":\"end\", \"status\":["
	for _, country := range countries {
		st := country.GetSummary()
		sum += st + ","
	}
	sum = sum[:len(sum)-1]
	sum += "]}"
	server.Ver.Lock()
	server.WriteMessage([]byte(sum))
	server.Ver.Unlock()
}
