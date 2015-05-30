package classification

import (
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/util"
	. "gopkg.in/check.v1"
	"math"
)

func (_ *Tests) TestWeakestLink(c *C) {
	tree := &Tree{
		leftChild: &Tree{
			leftChild: &Tree{
				counts: util.Histogram{4, 8},
			},
			rightChild: &Tree{
				counts: util.Histogram{8, 15},
			},
			counts: util.Histogram{12, 23},
		},
		rightChild: &Tree{
			leftChild: &Tree{
				counts: util.Histogram{4, 8},
			},
			rightChild: &Tree{
				counts: util.Histogram{8, 14},
			},
			counts: util.Histogram{12, 22},
		},
		counts: util.Histogram{24, 45},
	}

	ctx := &pruneCtx{
		lowestPenalty: math.Inf(1),
		pruneScore:    impurity.Gini,
	}
	ctx.findWeakestLink(tree, nil)

	c.Assert(ctx.bestPath, DeepEquals, []direction{left})
	found := ctx.lowestPenalty
	expected := impurity.Gini(util.Histogram{12, 23}) -
		impurity.Gini(util.Histogram{4, 8}) -
		impurity.Gini(util.Histogram{8, 15})
	c.Check(math.Abs(found-expected) <= 1e-6, Equals, true)
}

func (_ *Tests) TestCollapse(c *C) {
	tree1 := &Tree{
		counts: util.Histogram{0, 1},
	}
	for i := 0; i < 5; i++ {
		tree1 = &Tree{
			leftChild: &Tree{
				counts: util.Histogram{1, 0},
			},
			rightChild: tree1,
			counts:     util.Histogram{tree1.counts[0] + 1, tree1.counts[1]},
		}
	}

	// collapse the root
	tree2 := collapseSubtree(tree1, []direction{})
	c.Check(tree1.counts, DeepEquals, tree2.counts)
	c.Check(tree2.leftChild, IsNil)
	c.Check(tree2.rightChild, IsNil)

	// "collapse" a leaf
	tree3 := collapseSubtree(tree1, []direction{left})
	c.Check(tree1, Not(Equals), tree3)
	c.Check(tree1.rightChild, Equals, tree3.rightChild)
	c.Check(tree3.leftChild.counts, DeepEquals, tree1.leftChild.counts)
	c.Check(tree3.leftChild.leftChild, IsNil)
	c.Check(tree3.leftChild.rightChild, IsNil)

	// collapse an interior node
	path := []direction{right, right, right}
	tree4 := collapseSubtree(tree1, path)
	for i := 0; i < len(path); i++ {
		c.Check(tree1.counts, DeepEquals, tree4.counts)
		c.Check(tree1, Not(Equals), tree4)
		c.Check(tree1.leftChild, Equals, tree4.leftChild)
		tree1 = tree1.rightChild
		tree4 = tree4.rightChild
	}
	c.Check(tree1.counts, DeepEquals, tree4.counts)
	c.Check(tree1.leftChild, Not(IsNil))
	c.Check(tree1.rightChild, Not(IsNil))
	c.Check(tree4.leftChild, IsNil)
	c.Check(tree4.rightChild, IsNil)
}
