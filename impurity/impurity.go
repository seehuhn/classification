// impurity.go -
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

package impurity

import (
	"math"
)

type Function func([]int) float64

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

func sum(freq []int) int {
	res := 0
	for _, ni := range freq {
		res += ni
	}
	return res
}
