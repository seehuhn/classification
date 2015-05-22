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

// DevianceLoss computes -2 times the log-likelihood of getting the
// outcome `y` from the probability distribution with weigths `prob`.
func DevianceLoss(y int, prob []float64) float64 {
	return -2 * math.Log(prob[y])
}

func ZeroOneLoss(y int, prob []float64) float64 {
	var max float64
	for _, p := range prob {
		if p > max {
			max = p
		}
	}

	hit := 0
	count := 0
	for j, p := range prob {
		if p > max-1e-6 {
			count += 1
			if y == j {
				hit = 1
			}
		}
	}

	return 1.0 - float64(hit)/float64(count)
}

func OtherLoss(y int, prob []float64) float64 {
	q := 1 - prob[y]
	return q * q
}
