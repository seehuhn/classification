package tree

import (
	. "gopkg.in/check.v1"
)

func (*Tests) TestClasses(c *C) {
	tree := &Tree{
		LeftChild: &Tree{
			LeftChild: &Tree{
				Hist: []float64{1, 0, 0},
			},
			RightChild: &Tree{
				Hist: []float64{0, 1, 0},
			},
		},
		RightChild: &Tree{
			Hist: []float64{0, 0, 1},
		},
	}
	c.Check(tree.NumClasses(), Equals, 3)
}
