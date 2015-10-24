// xval.go - functions relating to cross-validation
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

package tree

import (
	"math/rand"
)

const xValSeed = 1769149487

func getXValSets(k, K, n int) (trainingSet, testSet []int) {
	if K < 2 {
		panic("need at least K=2 groups for cross-validation")
	}
	if n < K {
		panic("not enough samples for cross-validation")
	}

	a := k * n / K
	b := (k + 1) * n / K

	rng := rand.New(rand.NewSource(xValSeed))
	perm := rng.Perm(n)
	testSet = make([]int, b-a)
	copy(testSet, perm[a:b])
	trainingSet = append(perm[:a], perm[b:]...)

	return trainingSet, testSet
}
