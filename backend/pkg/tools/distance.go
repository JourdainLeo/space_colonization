package tools

import (
	plt "ai30-spatial/pkg/planet"
	"math"
	"sort"
)

func DistanceBetweenTwoStar(star1, star2 *plt.Star) float64 {
	p1X, p1Y, p1Z := star1.GetPosition()
	p2X, p2Y, p2Z := star2.GetPosition()
	return math.Pow(math.Pow(float64(p2X)-float64(p1X), 2)+math.Pow(float64(p2Y)-float64(p1Y), 2)+math.Pow(float64(p2Z)-float64(p1Z), 2), 0.5)
}

func DistanceBetweenStarAndCoordinates(star1 *plt.Star, posX, posY, posZ float64) float64 {
	p1X, p1Y, p1Z := star1.GetPosition()
	return math.Pow(math.Pow(float64(posX)-float64(p1X), 2)+math.Pow(float64(posY)-float64(p1Y), 2)+math.Pow(float64(posZ)-float64(p1Z), 2), 0.5)
}

func DistanceBetweenCoordinates(posX1, posY1, posZ1, posX2, posY2, posZ2 float64) float64 {
	return math.Pow(math.Pow(float64(posX2)-float64(posX1), 2)+math.Pow(float64(posY2)-float64(posY1), 2)+math.Pow(float64(posZ2)-float64(posZ1), 2), 0.5)
}

func SortPlanet(observed map[*plt.Star][]int, ord map[*plt.Star][]float64, strat string) []*plt.Star {

	ordered := make([]*plt.Star, 0)
	for key := range ord {
		ordered = append(ordered, key)
	}
	if strat == "quick" {
		sort.Slice(ordered, func(i, j int) bool {
			return observed[ordered[i]][0] > observed[ordered[j]][0]
		})
	} else if strat == "slow" {
		sort.Slice(ordered, func(i, j int) bool {
			return (ord[ordered[i]][0] / ord[ordered[i]][1]) > (ord[ordered[j]][0] / ord[ordered[j]][1])
		})

	}
	return ordered
}

func ShufflePlanet(planet map[*plt.Star][]int) []*plt.Star {
	shuffled := make([]*plt.Star, 0)
	for key := range planet {
		shuffled = append(shuffled, key)
	}
	return shuffled
}

func Contains(stars []*plt.Star, star *plt.Star) bool {
	for _, s := range stars {
		if s == star {
			return true
		}
	}
	return false
}

func ContainsMap(stars map[*plt.Star][]float64, star *plt.Star) bool {
	for key := range stars {
		if key == star {
			return true
		}
	}
	return false
}

func MinMapInt(m map[*plt.Star][]int) int {
	minVal := math.MaxInt64
	for _, val := range m {
		if val[0] < minVal {
			minVal = val[0]
		}
	}
	return minVal
}

func GetReaction(coef []int, strat string) []int {
	if strat == "faster" {
		return []int{coef[0] + 10, coef[1] + 10}
	} else if strat == "slower" {
		lo := coef[0] - 10
		hi := coef[0] - 10
		if lo <= 0 {
			lo = 5
		}
		if hi <= 0 {
			hi = 5
		}
		return []int{lo, hi}
	} else if strat == "same" {
		return coef
	} else if strat == "balanced_lo" {
		lo := coef[0] - 10
		if lo <= 0 {
			lo = 5
		}
		return []int{lo, coef[1] + 5}
	} else if strat == "balanced_hi" {
		hi := coef[0] - 10
		if hi <= 0 {
			hi = 10
		}
		return []int{coef[0] + 5, hi}
	}
	return coef
}

func GetPlanetName(planet []*plt.Star) string {
	name := make([]string, 0)
	for _, p := range planet {
		name = append(name, p.Name)
	}
	explode := "\"planets\": ["
	added := false
	for _, n := range name {
		explode += "\"" + n + "\", "
		added = true
	}
	if added {
		explode = explode[:len(explode)-2]
	}
	explode += "]"
	return explode
}

func GetRocketVelocity(planet *plt.Star) float64 {

	switch planet.Name {
	case "mercury":
		return 10
	case "venus":
		return 10
	case "mars":
		return 10
	case "jupiter":
		return 15
	case "saturn":
		return 25
	case "uranus":
		return 30
	case "neptune":
		return 35
	case "pluto":
		return 40
	}
	return 10
}
