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
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/stop"
	"github.com/seehuhn/classification/util"
	"math"
	"sort"
)

// Factory is a structure to store the parameters governing the
// growing and pruning of classification trees.  Any zero field values
// are interpreted as the corresponding values from the
// `DefaultFactory` structure.
type Factory struct {
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
var CART = &Factory{
	StopGrowth: stop.IfPureOrAtMost(5),
	SplitScore: impurity.Gini,
	PruneScore: impurity.MisclassificationError,
	XValLoss:   loss.ZeroOne,
	K:          10, // p.75
}

// DefaultFactory specifies the default parameters for constructing a
// tree; see the `Factory` documentation for the meaning of the
// individual fields.  The values given in `DefaultFactory` are used
// by the `FromData` function, and to replace zero values in a
// `Factory` structure when the `Factory.FromData` method is called.
var DefaultFactory = CART

// FromData constructs a new classification tree from training data.
func (b *Factory) FromData(data *data.Data) *Tree {
	n, p := data.X.Shape()
	if p > maxColumns {
		panic("too large p")
	}
	if len(data.Y) != n {
		panic("dimensions of `x` and `response` don't match")
	}
	if data.Rows != nil {
		n = len(data.Rows)
	}
	b = b.setDefaults()

	// step 1: generate the full tree
	tree := b.fullTree(data)
	if tree.IsLeaf() {
		return tree
	}

	// step 2: generate candidates for a pruned tree
	candidates, alpha := b.getCandidates(tree, data.NumClasses)
	loss := make([]float64, len(candidates))
	for k := 0; k < b.K; k++ {
		learnRows, testRows := getXValSets(k, b.K, n)

		// build the initial tree
		learnHist := util.GetHist(learnRows, data.NumClasses, data.Y, data.Weights)
		tree := b.getFullTree(data, learnRows, learnHist)

		// get all candidates for pruning the tree
		XVcandidates, XValpha := b.getCandidates(tree, data.NumClasses)

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
					prob := tree.EstimateClassProbabilities(data.X.Row(row))
					val := b.XValLoss(data.Y[row], prob)
					cumLoss += val
				}
				XVloss[i] = cumLoss
				XVlossDone[i] = true
			}
			loss[j] += XVloss[i]
		}
	}

	// step 3: return the optimal tree
	bestIdx := 0
	bestLoss := math.Inf(+1)
	for j := range candidates {
		if loss[j] <= bestLoss {
			bestIdx = j
			bestLoss = loss[j]
		}
	}
	return candidates[bestIdx]
}

func (b *Factory) setDefaults() *Factory {
	res := *b // make a copy
	if res.XValLoss == nil {
		res.XValLoss = DefaultFactory.XValLoss
	}
	if res.K == 0 {
		res.K = DefaultFactory.K
	}
	if res.StopGrowth == nil {
		res.StopGrowth = DefaultFactory.StopGrowth
	}
	if res.SplitScore == nil {
		res.SplitScore = DefaultFactory.SplitScore
	}
	if res.PruneScore == nil {
		res.PruneScore = DefaultFactory.PruneScore
	}
	return &res
}

func (b *Factory) fullTree(data *data.Data) *Tree {
	rows := data.GetRows()
	hist := util.GetHist(rows, data.NumClasses, data.Y, data.Weights)
	return b.getFullTree(data, rows, hist)
}

func (b *Factory) getFullTree(data *data.Data, rows []int, hist util.Histogram) *Tree {
	// TODO(voss): use a multi-threaded algorithm?

	if b.StopGrowth(hist) {
		return &Tree{
			Hist: hist,
		}
	}

	best := b.findBestSplit(data, rows, hist)

	return &Tree{
		Hist:       hist,
		LeftChild:  b.getFullTree(data, best.Left, best.LeftHist),
		RightChild: b.getFullTree(data, best.Right, best.RightHist),
		Column:     best.Col,
		Limit:      best.Limit,
	}
}

func (b *Factory) findBestSplit(data *data.Data, rows []int, hist util.Histogram) *searchResult {
	best := &searchResult{}
	first := true
	_, p := data.X.Shape()
	for col := 0; col < p; col++ {
		rows = copyIntSlice(rows)
		sort.Sort(&colSort{data.X, rows, col})

		leftHist := make(util.Histogram, len(hist))
		var rightHist = copyFloatSlice(hist)
		for i := 1; i < len(rows); i++ {
			yi := data.Y[rows[i-1]]
			leftHist[yi]++
			rightHist[yi]--

			left := data.X.At(rows[i-1], col)
			right := data.X.At(rows[i], col)
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
