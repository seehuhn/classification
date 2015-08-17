package classification

import (
	"bytes"
	"encoding"
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
	tree1.WriteBinary(w)
	r := bytes.NewReader(w.Bytes())
	tree2, err := TreeFromFile(r)
	c.Assert(err, Equals, nil)

	c.Check(tree1, DeepEquals, tree2)
}

func (*Tests) TestMarshalling(c *C) {
	tree1 := &Tree{
		LeftChild: &Tree{
			Hist: []int{0, 0, 10},
		},
		RightChild: &Tree{
			LeftChild: &Tree{
				Hist: []int{1, 0, 0},
			},
			RightChild: &Tree{
				Hist: []int{0, 1, 2},
			},
			Hist: []int{1, 1, 2},
		},
		Hist: []int{1, 1, 12},
	}

	data, err := tree1.MarshalBinary()
	c.Assert(err, Equals, nil)

	tree2 := &Tree{}
	err = tree2.UnmarshalBinary(data)
	c.Assert(err, Equals, nil)
	c.Assert(tree2, DeepEquals, tree1)
}

// compile time check: Tree implements encoding.BinaryMarshaler
var _ encoding.BinaryMarshaler = &Tree{}

// compile time check: Tree implements encoding.BinaryUnmarshaler
var _ encoding.BinaryUnmarshaler = &Tree{}
