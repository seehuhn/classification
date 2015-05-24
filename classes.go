// classes.go -
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

package classification

func probabilities(freq []int) []float64 {
	prob := make([]float64, len(freq))
	n := sum(freq)
	for i, ni := range freq {
		prob[i] = float64(ni) / float64(n)
	}
	return prob
}

func sum(freq []int) int {
	res := 0
	for _, ni := range freq {
		res += ni
	}
	return res
}
