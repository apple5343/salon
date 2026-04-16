package simulation

import (
	"math"
	"math/rand/v2"
)

func Poisson(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	for p > L {
		k++
		p *= rand.Float64()
	}
	return k - 1
}
