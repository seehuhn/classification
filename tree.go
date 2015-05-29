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
	"github.com/seehuhn/classification/util"
	"strings"
)

// Tree is the (opaque) data type to represent a classification tree.
type Tree struct {
	// fields used for internal nodes
	leftChild  *Tree
	rightChild *Tree
	column     int
	limit      float64

	// fields used for leaf nodes
	counts util.Histogram
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
			return t.counts.Probabilities()
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
// This function uses the `DefaultTreeBuilder` to construct a new
// classification tree.  The return values are the new tree and an
// estimate for the average value of the loss function (given by
// `DefaultTreeBuilder.XValLoss`).
func NewTree(x *Matrix, classes int, response []int) (*Tree, float64) {
	return DefaultTreeBuilder.NewTree(x, classes, response)
}
