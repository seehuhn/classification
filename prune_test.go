package classification

import (
	"github.com/seehuhn/classification/impurity"
	. "gopkg.in/check.v1"
	"math"
)

func (*Tests) TestInitialPrune(c *C) {
	tree1 := &Tree{
		Hist: []int{6, 8},
		LeftChild: &Tree{
			Hist: []int{3, 6},
			LeftChild: &Tree{
				Hist: []int{1, 2},
			},
			RightChild: &Tree{
				Hist: []int{2, 4},
			},
		},
		RightChild: &Tree{
			Hist: []int{3, 2},
			LeftChild: &Tree{
				Hist: []int{1, 1},
			},
			RightChild: &Tree{
				Hist: []int{2, 1},
			},
		},
	}
	b := &TreeBuilder{
		PruneScore: impurity.Gini,
	}
	b.initialPrune(tree1)
	tree2 := &Tree{
		Hist: []int{6, 8},
		LeftChild: &Tree{
			Hist: []int{3, 6},
		},
		RightChild: &Tree{
			Hist: []int{3, 2},
			LeftChild: &Tree{
				Hist: []int{1, 1},
			},
			RightChild: &Tree{
				Hist: []int{2, 1},
			},
		},
	}
	c.Assert(tree1, DeepEquals, tree2)
}

func (*Tests) TestWeakestLink(c *C) {
	tree := &Tree{
		Hist: []int{24, 45},
		LeftChild: &Tree{
			Hist: []int{12, 23},
			LeftChild: &Tree{
				Hist: []int{4, 8},
			},
			RightChild: &Tree{
				Hist: []int{8, 15},
			},
		},
		RightChild: &Tree{
			Hist: []int{12, 22},
			LeftChild: &Tree{
				Hist: []int{4, 8},
			},
			RightChild: &Tree{
				Hist: []int{8, 14},
			},
		},
	}

	ctx := &pruneCtx{
		lowestPenalty: math.Inf(1),
		pruneScore:    impurity.Gini,
	}
	ctx.findWeakestLink(tree, nil)

	c.Assert(ctx.bestPath, DeepEquals, []direction{left})
	found := ctx.lowestPenalty
	expected := impurity.Gini([]int{12, 23}) -
		impurity.Gini([]int{4, 8}) -
		impurity.Gini([]int{8, 15})
	c.Check(math.Abs(found-expected) <= 1e-6, Equals, true)
}

func (*Tests) TestCollapse(c *C) {
	tree1 := &Tree{
		Hist: []int{0, 1},
	}
	for i := 0; i < 5; i++ {
		tree1 = &Tree{
			LeftChild: &Tree{
				Hist: []int{1, 0},
			},
			RightChild: tree1,
			Hist:       []int{tree1.Hist[0] + 1, tree1.Hist[1]},
		}
	}

	// collapse the root
	tree2 := collapseSubtree(tree1, []direction{})
	c.Check(tree1.Hist, DeepEquals, tree2.Hist)
	c.Check(tree2.LeftChild, IsNil)
	c.Check(tree2.RightChild, IsNil)

	// "collapse" a leaf
	tree3 := collapseSubtree(tree1, []direction{left})
	c.Check(tree1, Not(Equals), tree3)
	c.Check(tree1.RightChild, Equals, tree3.RightChild)
	c.Check(tree3.LeftChild.Hist, DeepEquals, tree1.LeftChild.Hist)
	c.Check(tree3.LeftChild.LeftChild, IsNil)
	c.Check(tree3.LeftChild.RightChild, IsNil)

	// collapse an interior node
	path := []direction{right, right, right}
	tree4 := collapseSubtree(tree1, path)
	for i := 0; i < len(path); i++ {
		c.Check(tree1.Hist, DeepEquals, tree4.Hist)
		c.Check(tree1, Not(Equals), tree4)
		c.Check(tree1.LeftChild, Equals, tree4.LeftChild)
		tree1 = tree1.RightChild
		tree4 = tree4.RightChild
	}
	c.Check(tree1.Hist, DeepEquals, tree4.Hist)
	c.Check(tree1.LeftChild, Not(IsNil))
	c.Check(tree1.RightChild, Not(IsNil))
	c.Check(tree4.LeftChild, IsNil)
	c.Check(tree4.RightChild, IsNil)
}

func (*Tests) TestSelectCandidate(c *C) {
	alpha := []float64{0.0, 1.0, math.Inf(+1)}
	c.Check(selectCandidate(alpha, 0.0), Equals, 0)
	c.Check(selectCandidate(alpha, 0.5), Equals, 0)
	c.Check(selectCandidate(alpha, 2.0), Equals, 1)
	c.Check(selectCandidate(alpha, math.Inf(+1)), Equals, 1)
}
