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
	"math"
	"strings"
)

// Tree is the data type to represent nodes of a classification tree.
type Tree struct {
	// Hist gives the frequencies of the different reponses in the
	// trainings set, for this sub-tree.
	Hist util.Histogram

	// LeftChild points to the left subtree attached to this node.
	// For leaf nodes this is nil.
	LeftChild *Tree

	// RightChild points to the right subtree attached to this node.
	// For leaf nodes this is nil.
	RightChild *Tree

	// Column specifies which input variable this node splits the data
	// at.  This field is unused for leaf nodes.
	Column int

	// Limit specifies the critical value for the input variable given
	// by `Column`.  If the observed value is less than or equal to
	// `Limit`, the value corresponds to the left subtree, otherwise
	// the value corresponds to the right subtree.
	Limit float64
}

func (t *Tree) doFormat(indent int) []string {
	pfx := strings.Repeat("  ", indent)
	res := []string{}
	if t.LeftChild != nil {
		res = append(res, pfx+fmt.Sprintf("if x[%d] <= %g:", t.Column, t.Limit))
		res = append(res, t.LeftChild.doFormat(indent+1)...)
		res = append(res, pfx+"else:")
		res = append(res, t.RightChild.doFormat(indent+1)...)
	} else {
		res = append(res, pfx+fmt.Sprintf("%v", t.Hist))
	}

	return res
}

func (t *Tree) Format() string {
	return strings.Join(t.doFormat(0), "\n")
}

func (t *Tree) String() string {
	nodes := 0
	maxDepth := 0
	t.ForeachLeaf(func(_ util.Histogram, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		nodes++
	})
	return fmt.Sprintf("<classification tree, %d leaves, max depth %d>",
		nodes, maxDepth)
}

func (t *Tree) Lookup(x []float64) []float64 {
	for {
		if t.LeftChild == nil {
			return t.Hist.Probabilities()
		}
		if x[t.Column] <= t.Limit {
			t = t.LeftChild
		} else {
			t = t.RightChild
		}
	}
}

func (t *Tree) walkPostOrder(fn func(*Tree, int), depth int) {
	if t.LeftChild != nil {
		t.LeftChild.walkPostOrder(fn, depth+1)
		t.RightChild.walkPostOrder(fn, depth+1)
	}
	fn(t, depth)
}

func (t *Tree) ForeachLeaf(fn func(util.Histogram, int)) {
	t.foreachLeafRecursive(fn, 0)
}

func (t *Tree) foreachLeafRecursive(fn func(util.Histogram, int), depth int) {
	if t.LeftChild != nil {
		t.LeftChild.foreachLeafRecursive(fn, depth+1)
		t.RightChild.foreachLeafRecursive(fn, depth+1)
	} else {
		fn(t.Hist, depth)
	}
}

type RegionFunction func(a, b []float64, hist util.Histogram)

func (t *Tree) ForeachLeafRegion(p int, fn RegionFunction) {
	a := make([]float64, p)
	b := make([]float64, p)
	for i := 0; i < p; i++ {
		a[i] = math.Inf(-1)
		b[i] = math.Inf(+1)
	}
	t.foreachLeafRegionRecursive(a, b, fn)
}

func (t *Tree) foreachLeafRegionRecursive(a, b []float64, fn RegionFunction) {
	if t.LeftChild == nil {
		fn(a, b, t.Hist)
	} else {
		ai := a[t.Column]
		bi := b[t.Column]

		b[t.Column] = t.Limit
		t.LeftChild.foreachLeafRegionRecursive(a, b, fn)

		a[t.Column] = t.Limit
		b[t.Column] = bi
		t.RightChild.foreachLeafRegionRecursive(a, b, fn)

		a[t.Column] = ai
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
