package forest

import (
	"math"
	"math/rand"
	"sort"

	"seehuhn.de/go/classification/data"
	"seehuhn.de/go/classification/matrix"
)

func (f *RandomTree) findBestSplit(rng *rand.Rand, d *data.Data, hist data.Histogram) *searchResult {
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
	numColumns := f.NumColumns
	if numColumns == 0 {
		numColumns = int(math.Ceil(math.Sqrt(float64(d.NCol()))))
	}
	columns := subset(rng, numColumns, d.NCol())
	for _, col := range columns {
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

			leftScore := f.SplitScore(leftHist)
			rightScore := f.SplitScore(rightHist)
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

func copyIntSlice(src []int) []int {
	res := make([]int, len(src))
	copy(res, src)
	return res
}

func copyFloatSlice(src []float64) []float64 {
	res := make([]float64, len(src))
	copy(res, src)
	return res
}

// subset returns a random subset of {0, 1, ..., p-1} with m elements.
func subset(r *rand.Rand, m, p int) []int {
	if m > p {
		panic("invalid subset size")
	}

	// use reservoir sampling:
	res := make([]int, m)
	for i := 0; i < m; i++ {
		res[i] = i
	}
	for i := m; i < p; i++ {
		j := r.Intn(i + 1)
		if j < m {
			res[j] = i
		}
	}
	return res
}
