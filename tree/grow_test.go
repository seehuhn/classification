package tree

import (
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/matrix"
	. "gopkg.in/check.v1"
)

func (*Tests) TestFindBestSplit1(c *C) {
	n := 1000
	k := n / 3
	raw := make([]float64, n)
	response := make([]int, n)
	for i := range raw {
		raw[i] = float64(i)
		if i < k {
			response[i] = 1
		}
	}

	b := &Factory{
		SplitScore: impurity.Gini,
	}
	theData := &data.Data{
		NumClasses: 2,
		X:          matrix.NewFloat64(n, 1, 0, raw),
		Y:          response,
		Rows:       intRange(n),
	}
	total := make([]float64, theData.NumClasses)
	for _, row := range theData.Rows {
		yi := theData.Y[row]
		total[yi]++
	}
	best := b.findBestSplit(theData, theData.Rows, total)
	if len(best.Left) != k || len(best.Right) != n-k {
		c.Error("wrong split: expected", k, "got", len(best.Left))
	}
}

func (*Tests) TestFindBestSplit2(c *C) {
	//   0 1 2 3 4 5 6
	// -+-------------
	// 0|1 1 1 1 1 0 0
	// 1|1 1 1 1 1 0 0
	// 2|1 1 1 1 1 0 0
	n1 := 3
	n2 := 7
	k2 := 5
	raw := make([]float64, 2*n1*n2)
	response := make([]int, n1*n2)
	pos := 0
	for i1 := 0; i1 < n1; i1++ {
		for i2 := 0; i2 < n2; i2++ {
			raw[2*pos] = float64(i1)
			raw[2*pos+1] = float64(i2)
			if i2 < k2 {
				response[pos] = 1
			}
			pos++
		}
	}

	b := &Factory{
		SplitScore: impurity.Gini,
	}
	theData := &data.Data{
		NumClasses: 2,
		X:          matrix.NewFloat64(n1*n2, 2, 0, raw),
		Y:          response,
		Rows:       intRange(n1 * n2),
	}
	total := make([]float64, theData.NumClasses)
	for _, row := range theData.Rows {
		yi := response[row]
		total[yi]++
	}
	best := b.findBestSplit(theData, theData.Rows, total)
	if len(best.Left) != k2*n1 || len(best.Right) != (n2-k2)*n1 {
		c.Error("wrong split: expected ", k2*n1, " ", (n2-k2)*n1,
			" got ", len(best.Left), " ", len(best.Right))
	}
}

func (*Tests) TestFindBestSplit3(c *C) {
	n1 := 29
	k1 := 11
	n2 := 37
	k2 := 17
	raw := make([]float64, 2*n1*n2)
	response := make([]int, n1*n2)
	pos := 0
	for i1 := 0; i1 < n1; i1++ {
		for i2 := 0; i2 < n2; i2++ {
			raw[2*pos] = float64(i1)
			raw[2*pos+1] = float64(21)
			if i1 < k1 {
				response[pos] = 1
			} else if i2 < k2 {
				response[pos] = 2
			}
			pos++
		}
	}

	b := &Factory{
		SplitScore: impurity.Gini,
	}
	theData := &data.Data{
		NumClasses: 3,
		X:          matrix.NewFloat64(n1*n2, 2, 0, raw),
		Y:          response,
		Rows:       intRange(n1 * n2),
	}
	total := make([]float64, theData.NumClasses)
	for _, row := range theData.Rows {
		yi := theData.Y[row]
		total[yi]++
	}
	best := b.findBestSplit(theData, theData.Rows, total)
	if len(best.Left) != k1*n2 || len(best.Right) != (n1-k1)*n2 {
		c.Error("wrong split: expected", k1*n2, (n1-k1)*n2,
			"got", len(best.Left), len(best.Right))
	}
}
