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
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/util"
	"math"
	"strings"
)

// To prevent excessive memory use, the number of columns is limited
// to p <= maxColumns.
const maxColumns = 10000

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
	pfx := strings.Repeat("    ", indent)
	res := []string{}
	res = append(res, pfx+fmt.Sprintf("# %v", t.Hist))
	if !t.IsLeaf() {
		res = append(res, pfx+fmt.Sprintf("if x[%d] <= %g:", t.Column, t.Limit))
		res = append(res, t.LeftChild.doFormat(indent+1)...)
		res = append(res, pfx+"else:")
		res = append(res, t.RightChild.doFormat(indent+1)...)
	}

	return res
}

// Format returns a human readable, textual representation of the tree `t`.
func (t *Tree) Format() string {
	return strings.Join(t.doFormat(0), "\n")
}

// String returns a one-line summary description of the tree.  Use the
// `Format` method to get a textual representation of the full tree.
func (t *Tree) String() string {
	nodes := 0
	maxDepth := 0
	t.ForeachLeaf(func(_ util.Histogram, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		nodes++
	})
	tmpl := "<classification tree, %d leaves, max depth %d, representing %d samples>"
	return fmt.Sprintf(tmpl, nodes, maxDepth, t.Hist.Sum())
}

// Classes returns the number of classes of the response variable
// corresponding to the tree `t`.
func (t *Tree) Classes() int {
	p := len(t.Hist)
	if p > 0 {
		return p
	}
	return t.LeftChild.Classes()
}

// IsLeaf returns true if `t` is a terminal node and returns false if
// `t` has child nodes.
func (t *Tree) IsLeaf() bool {
	return t.LeftChild == nil
}

// lookup returns the terminal node corresponding to input `x`.
func (t *Tree) lookup(x []float64) *Tree {
	for {
		if t.IsLeaf() {
			return t
		}
		if x[t.Column] <= t.Limit {
			t = t.LeftChild
		} else {
			t = t.RightChild
		}
	}
}

// GetClassCounts returns the class counts for input `x`, as seen in
// the trainings data.
func (t *Tree) GetClassCounts(x []float64) util.Histogram {
	return t.lookup(x).Hist
}

// EstimateClassProbabilities returns the estimated class
// probabilities for input `x`.
func (t *Tree) EstimateClassProbabilities(x []float64) []float64 {
	return t.lookup(x).Hist.Probabilities()
}

// GuessClass tries to guess the class corresponding to input `x`.
func (t *Tree) GuessClass(x []float64) int {
	hist := t.lookup(x).Hist
	bestClass := -1
	bestCount := -1
	for class, count := range hist {
		if count > bestCount {
			bestCount = count
			bestClass = class
		}
	}
	return bestClass
}

func (t *Tree) walkPostOrder(fn func(*Tree, int), depth int) {
	if t.LeftChild != nil {
		t.LeftChild.walkPostOrder(fn, depth+1)
		t.RightChild.walkPostOrder(fn, depth+1)
	}
	fn(t, depth)
}

// ForeachLeaf calls the function `fn` once for each terminal node of
// the tree `t`.  The arguments to `fn` are the class counts for the
// samples corresponding to the node, and the depth of the node in the
// tree.
func (t *Tree) ForeachLeaf(fn func(hist util.Histogram, depth int)) {
	t.foreachLeafRecursive(0, fn)
}

func (t *Tree) foreachLeafRecursive(depth int, fn func(util.Histogram, int)) {
	if t.LeftChild != nil {
		t.LeftChild.foreachLeafRecursive(depth+1, fn)
		t.RightChild.foreachLeafRecursive(depth+1, fn)
	} else {
		fn(t.Hist, depth)
	}
}

// ForeachLeafRegion calls the function `fn` once for each terminal
// node in the tree `t`.  The arguments of `fn` are the rectangular
// region in feature space corresponding to the leaf node (`a` gives
// the minimal coordinate values, `b` gives the maximal coordinate
// values, negative infinities in `a` or positive infinities in `b`
// indicate unconstrained coordinates), the class counts for the
// samples corresponding to the node, as well as the depth of the node
// in the tree.
func (t *Tree) ForeachLeafRegion(
	fn func(a, b []float64, hist util.Histogram, depth int)) {
	p := t.Classes()
	a := make([]float64, p)
	b := make([]float64, p)
	for i := 0; i < p; i++ {
		a[i] = math.Inf(-1)
		b[i] = math.Inf(+1)
	}
	t.foreachLeafRegionRecursive(a, b, 0, fn)
}

func (t *Tree) foreachLeafRegionRecursive(a, b []float64, depth int,
	fn func(a, b []float64, hist util.Histogram, depth int)) {
	if t.IsLeaf() {
		fn(a, b, t.Hist, depth)
	} else {
		ai := a[t.Column]
		bi := b[t.Column]

		b[t.Column] = t.Limit
		t.LeftChild.foreachLeafRegionRecursive(a, b, depth+1, fn)

		a[t.Column] = t.Limit
		b[t.Column] = bi
		t.RightChild.foreachLeafRegionRecursive(a, b, depth+1, fn)

		a[t.Column] = ai
	}
}

// TreeFromTrainingsData constructs a new classification tree from a
// sample of trainings data.  The function uses the settings from
// `DefaultTreeBuilder`.
//
// The argument `classes` gives the number of classes in the response
// variable.  The rows of the matrix `x` are the observations from the
// trainings data set, and the corresponding classes are given as the
// entries of `y` (which must be in the range 0, 1, ..., classes-1).
//
// The return values are the new tree and a cross-validated estimate
// for the average value of the loss function (given by
// `DefaultTreeBuilder.XValLoss`).
func TreeFromTrainingsData(classes int, x *matrix.Float64, response []int) (*Tree, float64) {
	return DefaultTreeBuilder.TreeFromTrainingsData(classes, x, response)
}
