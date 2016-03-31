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
	"math"
	"sort"

	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/stop"
)

const xValSeed = 1769149487

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
	b = b.setDefaults()

	p := data.NCol()
	if p > maxColumns {
		panic("too large p")
	}

	// step 1: generate the full tree
	tree := b.fullTree(data)
	if tree.IsLeaf() {
		return tree
	}

	// step 2: generate candidates for a pruned tree
	candidates, alpha := b.getCandidates(tree)
	loss := make([]float64, len(candidates))
	for k := 0; k < b.K; k++ {
		xValSet := data.GetXValSet(xValSeed, b.K, k)

		// Build the initial tree using the training data.
		trainingData, _ := xValSet.TrainingData()
		tree := b.fullTree(trainingData)

		// Get all candidates for pruning the tree.
		XVcandidates, XValpha := b.getCandidates(tree)

		// Assess the expected loss of each candidate, using the test data.
		testData, _ := xValSet.TestData()
		testRows := testData.GetRows()
		XVloss := make([]float64, len(XVcandidates))
		XVlossDone := make([]bool, len(XVcandidates))
		for j := range candidates {
			a := math.Sqrt(alpha[j] * alpha[j+1])
			i := selectCandidate(XValpha, a)
			if !XVlossDone[i] {
				tree := XVcandidates[i]
				cumLoss := 0.0
				for _, row := range testRows {
					prob := tree.EstimateClassProbabilities(testData.X.Row(row))
					val := b.XValLoss(testData.Y[row], prob)
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
	hist := data.GetHist()
	return b.getFullTree(data, hist)
}

func (b *Factory) getFullTree(data *data.Data, hist data.Histogram) *Tree {
	if b.StopGrowth(hist) {
		return &Tree{
			Hist: hist,
		}
	}

	best := b.findBestSplit(data, hist)

	return &Tree{
		Hist:       hist,
		LeftChild:  b.getFullTree(best.Left, best.LeftHist),
		RightChild: b.getFullTree(best.Right, best.RightHist),
		Column:     best.Col,
		Limit:      best.Limit,
	}
}

func (b *Factory) findBestSplit(d *data.Data, hist data.Histogram) *searchResult {
	best := &searchResult{
		Left: &data.Data{
			NumClasses: d.NumClasses,
			X:          d.X,
			Y:          d.Y,
			Weights:    d.Weights,
		},
		Right: &data.Data{
			NumClasses: d.NumClasses,
			X:          d.X,
			Y:          d.Y,
			Weights:    d.Weights,
		},
	}
	first := true
	p := d.NCol()
	for col := 0; col < p; col++ {
		rows := copyIntSlice(d.GetRows())
		sort.Sort(&colSort{d.X, rows, col})

		leftHist := make(data.Histogram, len(hist))
		var rightHist = copyFloatSlice(hist)
		for i := 1; i < len(rows); i++ {
			yi := d.Y[rows[i-1]]
			leftHist[yi]++ // TODO(voss): use the weights here!
			rightHist[yi]--

			left := d.X.At(rows[i-1], col)
			right := d.X.At(rows[i], col)
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
				best.Left.Rows = rows[:i]
				best.Right.Rows = rows[i:]
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
	Left, Right         *data.Data
	LeftHist, RightHist data.Histogram
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
