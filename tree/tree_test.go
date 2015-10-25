package tree

import (
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type Tests struct{}

var _ = Suite(&Tests{})

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
