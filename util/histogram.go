// util.go - Auxiliary functions for github.com/seehuhn/classification
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

package util

// Histogram is the type used to represent class counts in a sample.
// The counts are stored as float64 values to allow for samples with
// non-integer weights.
type Histogram []float64

// GetHist counts how many instances of each class are seen in the
// given rows of the response data: `rows` specifies which entries of
// `y` to consider, `numClasses` gives the total number of possible
// classes, `y` gives the observed classes, and `w` gives the sample
// weights.  If `w` is nil, weight 1 is used for all samples.  The
// result is a Histogram of the (weighted) class counts.
func GetHist(rows []int, numClasses int, y []int, w []float64) Histogram {
	hist := make(Histogram, numClasses)
	if w == nil {
		for _, row := range rows {
			hist[y[row]]++
		}
	} else {
		for _, row := range rows {
			hist[y[row]] += w[row]
		}
	}
	return hist
}

// Sum returns the total number of samples corresponding to the
// histogram, obtained by adding up all entries of `hist`.
func (hist Histogram) Sum() float64 {
	res := 0.0
	for _, ni := range hist {
		res += ni
	}
	return res
}

// Probabilities returns an estimate of the class probabilities,
// obtained by normalising the entries of the histogram.
func (hist Histogram) Probabilities() Histogram {
	prob := make(Histogram, len(hist))
	total := hist.Sum()
	for i, ni := range hist {
		prob[i] = ni / total
	}
	return prob
}

// ArgMax returns the index of the histogram slot with the highest
// value.  In case of a draw, the lowest index involved is returned.
func (hist Histogram) ArgMax() int {
	bestIdx := 0
	bestVal := hist[0]
	for i := 1; i < len(hist); i++ {
		val := hist[i]
		if val > bestVal {
			bestIdx = i
			bestVal = val
		}
	}
	return bestIdx
}
