package classification

import (
	"math"
)

func Gini(freq []int) float64 {
	var res float64
	n := sum(freq)
	for _, ni := range freq {
		pi := float64(ni) / float64(n)
		res += pi * (1 - pi)
	}
	return res
}

func Entropy(freq []int) float64 {
	var res float64
	n := sum(freq)
	for _, ni := range freq {
		pi := float64(ni) / float64(n)
		if pi <= 1e-6 {
			continue
		}
		res -= pi * math.Log(pi)
	}
	return res
}

func MisclassificationError(freq []int) float64 {
	n := sum(freq)
	max := 0
	for _, ni := range freq {
		if ni > max {
			max = ni
		}
	}
	return float64(n-max) / float64(n)
}
