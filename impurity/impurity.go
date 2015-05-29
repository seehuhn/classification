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
	"github.com/seehuhn/classification/util"
	"math"
)

// Function is the type of functions used to compute impurity measures
// for a set of classification results.  The argument of a Function is
// a histogram, giving the counts for the different classes.  The
// output should be 0, if only one value is represented in the
// histogram, and positive otherwise.
//
// The output of a Function should scale linearly with the input
// vector, i.e. if the all entries of the input vector are doubled,
// the output value should double, too.
type Function func(util.Histogram) float64

func Gini(freq util.Histogram) float64 {
	var res float64
	n := float64(freq.Sum())
	for _, ni := range freq {
		floatNi := float64(ni)
		pi := floatNi / n
		res += floatNi * (1 - pi)
	}
	return res
}

func Entropy(freq util.Histogram) float64 {
	var res float64
	n := float64(freq.Sum())
	for _, ni := range freq {
		floatNi := float64(ni)
		pi := floatNi / n
		if pi <= 1e-6 {
			continue
		}
		res -= floatNi * math.Log(pi)
	}
	return res
}

func MisclassificationError(freq util.Histogram) float64 {
	n := freq.Sum()
	max := 0
	for _, ni := range freq {
		if ni > max {
			max = ni
		}
	}
	return float64(n - max)
}
