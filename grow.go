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

type xBuilder struct {
	TreeBuilder
	x        *Matrix
	classes  int
	response []int
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
