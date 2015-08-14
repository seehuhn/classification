package classification

import (
	"bytes"
	. "gopkg.in/check.v1"
)

func (*Tests) TestBinaryFormat(c *C) {
	tree1 := &Tree{
		LeftChild: &Tree{
			LeftChild: &Tree{
				Hist: []int{1, 0, 0},
			},
			RightChild: &Tree{
				Hist: []int{0, 1, 0},
			},
			Hist: []int{1, 1, 0},
		},
		RightChild: &Tree{
			Hist: []int{0, 0, 1},
		},
		Hist: []int{1, 1, 1},
	}

	w := &bytes.Buffer{}
	tree1.WriteTo(w)
	r := bytes.NewReader(w.Bytes())
	tree2, err := TreeFromFile(r)
	c.Assert(err, Equals, nil)

	c.Check(tree1, DeepEquals, tree2)
}
