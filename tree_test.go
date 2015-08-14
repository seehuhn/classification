package classification

import (
	. "gopkg.in/check.v1"
)

func (*Tests) TestClasses(c *C) {
	tree := &Tree{
		LeftChild: &Tree{
			LeftChild: &Tree{
				Hist: []int{1, 0, 0},
			},
			RightChild: &Tree{
				Hist: []int{0, 1, 0},
			},
		},
		RightChild: &Tree{
			Hist: []int{0, 0, 1},
		},
	}
	c.Check(tree.Classes(), Equals, 3)
}
