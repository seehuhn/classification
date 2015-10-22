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

package tree

import (
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/stop"
	"github.com/seehuhn/classification/util"
	"math"
	"sort"
)

// Builder is a structure to store the parameters governing the
// growing and pruning of classification trees.  Any zero field values
// are interpreted as the corresponding values from the
// `DefaultBuilder` structure.
type Builder struct {
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

// CART specifies the parameters for constructing a tree as suggested
// in the book "Classification and Regression Trees" by Breiman et
// al. (Chapman & Hall CRC, 1984).
var CART = &Builder{
	StopGrowth: stop.IfPureOrAtMost(5),
	SplitScore: impurity.Gini,
	PruneScore: impurity.MisclassificationError,
	XValLoss:   loss.ZeroOne,
	K:          10, // p.75
}

// DefaultBuilder specifies the default parameters for
// constructing a tree; see the `Builder` documentation for the
// meaning of the individual fields.  The values given in
// `DefaultBuilder` are used by the `NewFromTrainingData`
// function, and to replace zero values in a `Builder` structure
// when the `Builder.NewFromTrainingData` method is called.
var DefaultBuilder = CART

func (b *Builder) setDefaults() {
	if b.XValLoss == nil {
		b.XValLoss = DefaultBuilder.XValLoss
	}
	if b.K == 0 {
		b.K = DefaultBuilder.K
	}
	if b.StopGrowth == nil {
		b.StopGrowth = DefaultBuilder.StopGrowth
	}
	if b.SplitScore == nil {
		b.SplitScore = DefaultBuilder.SplitScore
	}
	if b.PruneScore == nil {
		b.PruneScore = DefaultBuilder.PruneScore
	}
}

func (b *Builder) fullTree(x *matrix.Float64, classes int, response []int, w []float64) *Tree {
	b.setDefaults()
	rows := intRange(len(response))
	hist := util.GetHist(rows, classes, response, w)
	xb := &xBuilder{*b, x, classes, response}
	return xb.getFullTree(rows, hist)
}

// NewFromTrainingData constructs a new classification tree from
// training data.
//
// K-fold crossvalidation is used to find the optimal pruning
// parameter.
//
// The return values are the new tree and an estimate for the average
// value of the loss function (given by `b.XValLoss`).
func (b *Builder) NewFromTrainingData(classes int, x *matrix.Float64,
	response []int, w []float64) (*Tree, float64) {
	b.setDefaults()

	n, p := x.Shape()
	if p > maxColumns {
		panic("too large p")
	}
	if len(response) != n {
		panic("dimensions of `x` and `response` don't match")
	}

	// generate the full tree
	tree := b.fullTree(x, classes, response, w)
	candidates, alpha := b.getCandidates(tree, classes)
	loss := make([]float64, len(candidates))

	for k := 0; k < b.K; k++ {
		learnRows, testRows := getXValSets(k, b.K, n)

		// build the initial tree
		learnHist := util.GetHist(learnRows, classes, response, w)
		xb := &xBuilder{*b, x, classes, response}
		tree := xb.getFullTree(learnRows, learnHist)

		// get all candidates for pruning the tree
		XVcandidates, XValpha := b.getCandidates(tree, classes)

		// for each alpha, assess the expected loss using the test set
		XVloss := make([]float64, len(XVcandidates))
		XVlossDone := make([]bool, len(XVcandidates))
		for j := range candidates {
			a := math.Sqrt(alpha[j] * alpha[j+1])
			i := selectCandidate(XValpha, a)
			if !XVlossDone[i] {
				tree := XVcandidates[i]
				cumLoss := 0.0
				for _, row := range testRows {
					prob := tree.EstimateClassProbabilities(x.Row(row))
					val := b.XValLoss(response[row], prob)
					cumLoss += val
				}
				XVloss[i] = cumLoss
				XVlossDone[i] = true
			}
			loss[j] += XVloss[i]
		}
	}

	// find the optimal tree
	bestIdx := 0
	bestLoss := math.Inf(+1)
	for j := range candidates {
		if loss[j] <= bestLoss {
			bestIdx = j
			bestLoss = loss[j]
		}
	}

	return candidates[bestIdx], bestLoss / float64(len(response))
}

type xBuilder struct {
	Builder
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
		var rightHist = copyFloatSlice(hist)
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
				best.LeftHist = copyFloatSlice(leftHist)
				best.RightHist = copyFloatSlice(rightHist)
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
