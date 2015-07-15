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

type Histogram []int

func GetHist(rows []int, classes int, y []int) Histogram {
	hist := make([]int, classes)
	for _, row := range rows {
		hist[y[row]]++
	}
	return hist
}

func (freq Histogram) Sum() int {
	res := 0
	for _, ni := range freq {
		res += ni
	}
	return res
}

func (freq Histogram) Probabilities() []float64 {
	prob := make([]float64, len(freq))
	n := freq.Sum()
	for i, ni := range freq {
		prob[i] = float64(ni) / float64(n)
	}
	return prob
}
