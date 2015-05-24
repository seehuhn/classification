// loss.go -
// Copyright (C) 2015  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
