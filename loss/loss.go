// Package loss implements loss functions for the assessment of model fit.
package loss

import (
	"math"
)

type Function func(int, []float64) float64

// Deviance computes -2 times the log-likelihood of getting the
// outcome `y` from the probability distribution with weigths `prob`.
func Deviance(y int, prob []float64) float64 {
	return -2 * math.Log(prob[y])
}

func ZeroOne(y int, prob []float64) float64 {
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

func Other(y int, prob []float64) float64 {
	q := 1 - prob[y]
	return q * q
}
