package classification

import (
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/util"
	"sort"
)

type TreeBuilder struct {
	XValLoss loss.Function
	K        int

	StopGrowth StopFunction
	SplitScore impurity.Function
	PruneScore impurity.Function
}

var DefaultTreeBuilder = &TreeBuilder{
	XValLoss:   loss.ZeroOne,
	K:          5,
	StopGrowth: StopIfAtMost(1),
	SplitScore: impurity.Gini,
	PruneScore: impurity.MisclassificationError,
}

func (b *TreeBuilder) setDefaults() {
	if b.XValLoss == nil {
		b.XValLoss = DefaultTreeBuilder.XValLoss
	}
	if b.K <= 0 {
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

// NewTree constructs a new classification tree.
//
// K-fold crossvalidation is used to find the optimal pruning
// parameter.
//
// The return values are the new tree and an estimate for the average
// value of the loss function (given by `b.XValLoss`).
func (b *TreeBuilder) NewTree(x *Matrix, classes int, response []int) (*Tree, float64) {
	b.setDefaults()

	n := len(response)

	alphaSteps := 50
	alpha := make([]float64, alphaSteps)
	alpha[0] = -1 // ask tryTrees to determine the range of alpha

	mean := make([]float64, len(alpha))

	for k := 0; k < b.K; k++ {
		learnRows, testRows := getXValSets(k, b.K, n)

		// build the initial tree
		learnHist := util.GetHist(learnRows, classes, response)
		xb := &xBuilder{*b, x, classes, response}
		tree := xb.getFullTree(learnRows, learnHist)

		// get all candidates for pruning the tree
		candidates := b.prunedTrees(tree, classes)

		// get the optimally pruned tree for every alpha
		trees := b.tryTrees(candidates, alpha)

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
	hist := util.GetHist(rows, classes, response)
	xb := &xBuilder{*b, x, classes, response}
	tree := xb.getFullTree(rows, hist)

	candidates := b.prunedTrees(tree, classes)
	tree = b.tryTrees(candidates, []float64{bestAlpha})[0]

	return tree, bestExpectedLoss
}

type xBuilder struct {
	TreeBuilder
	x        *Matrix
	classes  int
	response []int
}

func (b *xBuilder) getFullTree(rows []int, hist util.Histogram) *Tree {
	if b.StopGrowth(hist) {
		return &Tree{
			counts: hist,
		}
	}

	best := b.findBestSplit(rows, hist)

	return &Tree{
		leftChild:  b.getFullTree(best.Left, best.LeftHist),
		rightChild: b.getFullTree(best.Right, best.RightHist),
		column:     best.Col,
		limit:      best.Limit,
		counts:     hist,
	}
}

func (b *xBuilder) findBestSplit(rows []int, hist util.Histogram) *searchResult {
	best := &searchResult{}
	first := true
	for col := 0; col < b.x.p; col++ {
		rows = copyIntSlice(rows)
		sort.Sort(&colSort{b.x, rows, col})

		leftHist := make(util.Histogram, len(hist))
		var rightHist util.Histogram = copyIntSlice(hist)
		for i := 1; i < len(rows); i++ {
			yi := b.response[rows[i-1]]
			leftHist[yi]++
			rightHist[yi]--
			leftScore := b.SplitScore(leftHist)
			rightScore := b.SplitScore(rightHist)
			// TODO(voss): check that the score is computed correctly
			score := (leftScore + rightScore) / float64(len(rows))
			if first || score < best.Score {
				best.Col = col
				best.Limit = (b.x.At(rows[i-1], col) + b.x.At(rows[i], col)) / 2
				best.Left = rows[:i]
				best.Right = rows[i:]
				best.LeftHist = copyIntSlice(leftHist)
				best.RightHist = copyIntSlice(rightHist)
				best.Score = score
				first = false
			}
		}
		if rightHist.Sum() != 1 { // only rows[len(rows)-1] should be there
			panic("wrong histogram passed to findBestSplit")
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
	x    *Matrix
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
