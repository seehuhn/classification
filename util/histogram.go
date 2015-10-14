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
type Histogram []int

// GetHist counts how many instances of each class are seen in the
// given rows of the response data: `rows` specifies which entries of
// `y` to consider, `classes` gives the total number of possible
// classes, `y` gives the observed classes.  The result is a Histogram
// of the class counts.
func GetHist(rows []int, classes int, y []int) Histogram {
	hist := make([]int, classes)
	for _, row := range rows {
		hist[y[row]]++
	}
	return hist
}

// Sum returns the total number of samples corresponding to the
// histogram, obtained by adding up all entries of `hist`.
func (hist Histogram) Sum() int {
	res := 0
	for _, ni := range hist {
		res += ni
	}
	return res
}

// Probabilities returns an estimate of the class probabilities,
// obtained by normalising the entries of the histogram.
func (hist Histogram) Probabilities() []float64 {
	prob := make([]float64, len(hist))
	n := hist.Sum()
	for i, ni := range hist {
		prob[i] = float64(ni) / float64(n)
	}
	return prob
}
