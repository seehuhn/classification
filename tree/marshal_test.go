package tree

import (
	"bytes"
	"encoding"

	. "gopkg.in/check.v1"
)

func (*Tests) TestBinaryFormat(c *C) {
	tree1 := &Tree{
		LeftChild: &Tree{
			LeftChild: &Tree{
				Hist: []float64{1, 0, 0},
			},
			RightChild: &Tree{
				Hist: []float64{0, 1, 0},
			},
			Hist: []float64{1, 1, 0},
		},
		RightChild: &Tree{
			Hist: []float64{0, 0, 1},
		},
		Hist: []float64{1, 1, 1},
	}

	w := &bytes.Buffer{}
	tree1.WriteBinary(w)
	r := bytes.NewReader(w.Bytes())
	tree2, err := FromFile(r)
	c.Assert(err, Equals, nil)

	c.Check(tree1, DeepEquals, tree2)
}

func (*Tests) TestMarshalling(c *C) {
	tree1 := &Tree{
		LeftChild: &Tree{
			Hist: []float64{0, 0, 10},
		},
		RightChild: &Tree{
			LeftChild: &Tree{
				Hist: []float64{1, 0, 0},
			},
			RightChild: &Tree{
				Hist: []float64{0, 1, 2},
			},
			Hist: []float64{1, 1, 2},
		},
		Hist: []float64{1, 1, 12},
	}

	data, err := tree1.MarshalBinary()
	c.Assert(err, Equals, nil)

	tree2 := &Tree{}
	err = tree2.UnmarshalBinary(data)
	c.Assert(err, Equals, nil)
	c.Assert(tree2, DeepEquals, tree1)
}

func (*Tests) TestFuzzerCrash1(c *C) {
	data := []byte("JVCT\x01\xb4\x99ѿ\x02\x01\x01\xd1\xea\xf6\x18\xfd?\x00d" +
		"\x00")
	tree := &Tree{}
	err := tree.UnmarshalBinary(data)
	c.Assert(err, Equals, ErrTreeEncoding)
}

// compile time check: Tree implements encoding.BinaryMarshaler
var _ encoding.BinaryMarshaler = &Tree{}

// compile time check: Tree implements encoding.BinaryUnmarshaler
var _ encoding.BinaryUnmarshaler = &Tree{}
