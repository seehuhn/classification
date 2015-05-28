// crossvalidation.go -
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

func getXValSets(k, K, n int) (trainingSet, testSet []int) {
	trainingSetSize := n * (K - 1) / K
	trainingSet = make([]int, 0, trainingSetSize+1)
	testSet = make([]int, 0, n-trainingSetSize)
	for i := 0; i < n; i++ {
		if i%K == k {
			testSet = append(testSet, i)
		} else {
			trainingSet = append(trainingSet, i)
		}
	}
	return
}
