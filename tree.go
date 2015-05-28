// tree.go -
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

import (
	"fmt"
	"strings"
)

type StopFunction func([]int) bool

type Tree struct {
	// fields used for internal nodes
	leftChild  *Tree
	rightChild *Tree
	column     int
	limit      float64

	// fields used for leaf nodes
	counts []int
}

func (t *Tree) doFormat(indent int) []string {
	pfx := strings.Repeat("  ", indent)
	res := []string{}
	if t.counts != nil {
		res = append(res, pfx+fmt.Sprintf("%v", t.counts))
	}
	if t.leftChild != nil {
		res = append(res, pfx+fmt.Sprintf("if x[%d] <= %g:", t.column, t.limit))
		res = append(res, t.leftChild.doFormat(indent+1)...)
		res = append(res, pfx+"else:")
		res = append(res, t.rightChild.doFormat(indent+1)...)
	}
	return res
}

func (t *Tree) Format() string {
	return strings.Join(t.doFormat(0), "\n")
}

func (t *Tree) String() string {
	nodes := 0
	maxDepth := 0
	t.foreachLeaf(func(_ *Tree, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		nodes++
	}, 0)
	return fmt.Sprintf("<classification tree, %d leaves, max depth %d>",
		nodes, maxDepth)
}

func (t *Tree) Lookup(x []float64) []float64 {
	for {
		if t.leftChild == nil {
			return probabilities(t.counts)
		}
		if x[t.column] <= t.limit {
			t = t.leftChild
		} else {
			t = t.rightChild
		}
	}
}

func (t *Tree) walkPostOrder(fn func(*Tree, int), depth int) {
	if t.leftChild != nil {
		t.leftChild.walkPostOrder(fn, depth+1)
		t.rightChild.walkPostOrder(fn, depth+1)
	}
	fn(t, depth)
}

func (t *Tree) foreachLeaf(fn func(*Tree, int), depth int) {
	if t.leftChild != nil {
		t.leftChild.foreachLeaf(fn, depth+1)
		t.rightChild.foreachLeaf(fn, depth+1)
	} else {
		fn(t, depth)
	}
}

// NewTree constructs a new classification tree.
//
// K-fold crossvalidation is used to find the optimal pruning
// parameter.
func (b *TreeBuilder) NewTree(x *Matrix, classes int, response []int) (*Tree, float64) {
	n := len(response)

	alphaSteps := 50
	alpha := make([]float64, alphaSteps)
	alpha[0] = -1 // ask tryTrees to determine the range of alpha

	mean := make([]float64, len(alpha))

	for k := 0; k < b.K; k++ {
		learnRows, testRows := getXValSets(k, b.K, n)

		trees := b.tryTrees(x, classes, response, learnRows, alpha)

		cache := make(map[*Tree]float64)
		for l, tree := range trees {
			var cumLoss float64
			if loss, ok := cache[tree]; ok {
				cumLoss = loss
			} else {
				for _, row := range testRows {
					prob := tree.Lookup(x.Row(row))
					val := b.XValLoss(response[row], prob)
					cumLoss += val
				}
				cache[tree] = cumLoss
			}
			mean[l] += cumLoss
		}
	}
	for l := range alpha {
		mean[l] /= float64(n)
	}

	var bestAlpha float64
	var bestExpectedLoss float64
	for l, a := range alpha {
		if l == 0 || mean[l] < bestExpectedLoss {
			bestAlpha = a
			bestExpectedLoss = mean[l]
		}
	}

	rows := intRange(len(response))
	tree := b.tryTrees(x, classes, response, rows, []float64{bestAlpha})[0]

	return tree, bestExpectedLoss
}
