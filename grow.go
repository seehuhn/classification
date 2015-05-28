package classification

import (
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
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
	XValLoss: loss.ZeroOne,
	K:        5,
	StopGrowth: func(y []int) bool {
		return len(y) <= 1
	},
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

type xBuilder struct {
	TreeBuilder
	x        *Matrix
	classes  int
	response []int
}

func (b *xBuilder) build(rows []int) *Tree {
	y := make([]int, len(rows))
	hist := make([]int, b.classes)
	for i, row := range rows {
		yi := b.response[row]
		y[i] = yi
		hist[yi]++
	}

	if b.StopGrowth(y) {
		return &Tree{
			counts: hist,
		}
	}

	best := b.findBestSplit(rows, hist)

	return &Tree{
		leftChild:  b.build(best.Left),
		rightChild: b.build(best.Right),
		column:     best.Col,
		limit:      best.Limit,
	}
}

type searchResult struct {
	Col         int
	Limit       float64
	Left, Right []int
	Score       float64
}

func (b *xBuilder) findBestSplit(rows []int, hist []int) *searchResult {
	best := &searchResult{}
	first := true
	for col := 0; col < b.x.p; col++ {
		sort.Sort(&colSort{b.x, rows, col})

		leftFreq := make([]int, len(hist))
		rightFreq := copyIntSlice(hist)
		for i := 1; i < len(rows); i++ {
			limit := (b.x.At(rows[i-1], col) + b.x.At(rows[i], col)) / 2
			yi := b.response[rows[i-1]]
			leftFreq[yi]++
			rightFreq[yi]--
			leftScore := b.SplitScore(leftFreq)
			rightScore := b.SplitScore(rightFreq)
			p := float64(i) / float64(len(rows))
			score := leftScore*p + rightScore*(1-p)
			if first || score < best.Score {
				best.Col = col
				best.Limit = limit
				best.Left = copyIntSlice(rows[:i])
				best.Right = copyIntSlice(rows[i:])
				best.Score = score
				first = false
			}
		}
	}
	return best
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
