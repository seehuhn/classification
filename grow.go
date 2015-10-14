// grow.go - functions to grow a classification tree from data
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
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/stop"
	"github.com/seehuhn/classification/util"
	"sort"
)

// TreeBuilder is a structure to store the parameters governing the
// growing and pruning of classification trees.  Any zero field values
// are interpreted as the corresponding values from the
// `DefaultTreeBuilder` structure.
type TreeBuilder struct {
	// StopGrowth decides for every node of the initial tree whether a
	// further split is considered.  The default is to stop splitting
	// nodes once the node only contains a single type of
	// observations.
	StopGrowth stop.Function

	// SplitScore is the impurity function used when growing the
	// initial tree.  Splits are (greedily) chosen to minimise the
	// total `SplitScore`.  The default is to use Gini impurity here.
	SplitScore impurity.Function

	// PruneScore is the impurity function ("cost") used for
	// cost-complexity pruning.  Pruning aims to reduce the size
	// ("complexity") of the tree, while keeping the `PruneScore`
	// small.  The default is to use the misclassification error rate.
	PruneScore impurity.Function

	// XValLoss is used to guide the balance between the cost (as
	// measured by `PruneScore`) and complexity (i.e. the final size
	// of the tree).  Cost and complexity are balanced to minimise
	// XValLoss.  The default is to use zero-one loss.
	XValLoss loss.Function

	// The number of groups to use in cross-validation when estimating
	// the expected loss.  The default is to use 5 groups.
	K int
}

// DefaultTreeBuilder specifies the default parameters for
// constructing a tree; see the `TreeBuilder` documentation for the
// meaning of the individual fields.  The values given in
// `DefaultTreeBuilder` are used by the `TreeFromTrainingsData`
// function, and to replace zero values in a `TreeBuilder` structure
// when the `TreeBuilder.TreeFromTrainingsData` method is called.
var DefaultTreeBuilder = &TreeBuilder{
	StopGrowth: stop.IfHomogeneous,
	SplitScore: impurity.Gini,
	PruneScore: impurity.MisclassificationError,
	XValLoss:   loss.ZeroOne,
	K:          5,
}

func (b *TreeBuilder) setDefaults() {
	if b.XValLoss == nil {
		b.XValLoss = DefaultTreeBuilder.XValLoss
	}
	if b.K == 0 {
		b.K = DefaultTreeBuilder.K
	}
	if b.StopGrowth == nil {
		b.StopGrowth = DefaultTreeBuilder.StopGrowth
	}
	if b.SplitScore == nil {
		b.SplitScore = DefaultTreeBuilder.SplitScore
	}
	if b.PruneScore == nil {
		b.PruneScore = DefaultTreeBuilder.PruneScore
	}
}

func (b *TreeBuilder) fullTree(x *matrix.Float64, classes int, response []int) *Tree {
	b.setDefaults()
	rows := intRange(len(response))
	hist := util.GetHist(rows, classes, response)
	xb := &xBuilder{*b, x, classes, response}
	return xb.getFullTree(rows, hist)
}

// TreeFromTrainingsData constructs a new classification tree from
// trainings data.
//
// K-fold crossvalidation is used to find the optimal pruning
// parameter.
//
// The return values are the new tree and an estimate for the average
// value of the loss function (given by `b.XValLoss`).
func (b *TreeBuilder) TreeFromTrainingsData(classes int, x *matrix.Float64,
	response []int) (*Tree, float64) {
	b.setDefaults()

	n, p := x.Shape()
	if p > maxColumns {
		panic("too large p")
	}
	if len(response) != n {
		panic("dimensions of `x` and `response` don't match")
	}

	tune := &tuneProfile{}
	for k := 0; k < b.K; k++ {
		learnRows, testRows := getXValSets(k, b.K, n)

		// build the initial tree
		learnHist := util.GetHist(learnRows, classes, response)
		xb := &xBuilder{*b, x, classes, response}
		tree := xb.getFullTree(learnRows, learnHist)

		// get all candidates for pruning the tree
		b.initialPrune(tree)
		candidates := b.getCandidates(tree, classes)
		breaks, candidates := b.filterCandidates(candidates)

		// for each candidate, assess the expected loss using the test set
		losses := make([]float64, len(candidates))
		for i, tree := range candidates {
			cumLoss := 0.0
			for _, row := range testRows {
				prob := tree.EstimateClassProbabilities(x.Row(row))
				val := b.XValLoss(response[row], prob)
				cumLoss += val
			}
			losses[i] = cumLoss
		}

		tune.Add(breaks, losses)
	}

	// generate the optimal tree
	bestAlpha, bestExpectedLoss := tune.Minimum()

	tree := b.fullTree(x, classes, response)
	b.initialPrune(tree)
	candidates := b.getCandidates(tree, classes)
	tree = b.selectCandidate(candidates, bestAlpha)

	return tree, bestExpectedLoss / float64(len(response))
}

type xBuilder struct {
	TreeBuilder
	x        *matrix.Float64
	classes  int
	response []int
}

// potential plan to reduce the amount of sorting required
//
// 1. sort rows by col j: i0, i1, i2, ..., in only once
// 2. split rows: i0, ..., ik | i{k+1}, ..., in as before
// 3. for every other row:
//    - old sort order is: i0', i1', ..., in'
//    - after the split, the order stays the same, but elements are
//      sorted into two groups.

func (b *xBuilder) getFullTree(rows []int, hist util.Histogram) *Tree {
	// TODO(voss): use a multi-threaded algorithm?
	if b.StopGrowth(hist) {
		return &Tree{
			Hist: hist,
		}
	}

	best := b.findBestSplit(rows, hist)

	return &Tree{
		Hist:       hist,
		LeftChild:  b.getFullTree(best.Left, best.LeftHist),
		RightChild: b.getFullTree(best.Right, best.RightHist),
		Column:     best.Col,
		Limit:      best.Limit,
	}
}

func (b *xBuilder) findBestSplit(rows []int, hist util.Histogram) *searchResult {
	best := &searchResult{}
	first := true
	_, p := b.x.Shape()
	for col := 0; col < p; col++ {
		rows = copyIntSlice(rows)
		sort.Sort(&colSort{b.x, rows, col})

		leftHist := make(util.Histogram, len(hist))
		var rightHist = copyIntSlice(hist)
		for i := 1; i < len(rows); i++ {
			yi := b.response[rows[i-1]]
			leftHist[yi]++
			rightHist[yi]--

			left := b.x.At(rows[i-1], col)
			right := b.x.At(rows[i], col)
			if !(left < right) {
				continue
			}
			limit := (left + right) / 2

			leftScore := b.SplitScore(leftHist)
			rightScore := b.SplitScore(rightHist)
			score := leftScore + rightScore

			if first || score < best.Score {
				best.Col = col
				best.Limit = limit
				best.Left = rows[:i]
				best.Right = rows[i:]
				best.LeftHist = copyIntSlice(leftHist)
				best.RightHist = copyIntSlice(rightHist)
				best.Score = score
				first = false
			}
		}
	}
	return best
}

type searchResult struct {
	Col                 int
	Limit               float64
	Left, Right         []int
	LeftHist, RightHist util.Histogram
	Score               float64
}

type colSort struct {
	x    *matrix.Float64
	rows []int
	col  int
}

func (c *colSort) Len() int { return len(c.rows) }
func (c *colSort) Less(i, j int) bool {
	return c.x.At(c.rows[i], c.col) < c.x.At(c.rows[j], c.col)
}
func (c *colSort) Swap(i, j int) {
	c.rows[i], c.rows[j] = c.rows[j], c.rows[i]
}
