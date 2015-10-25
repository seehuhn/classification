package tree

import (
	"github.com/seehuhn/classification/impurity"
	"math"
)

const epsilon = 1e-6

// The initialPrune method modifies the given tree by recursively
// collapsing all leaves where the impurity `b.PruneScore` is not
// increased in the process.  The return value is the total impurity
// value of the pruned tree.
func (b *Factory) initialPrune(tree *Tree) float64 {
	thisScore := b.PruneScore(tree.Hist)
	if tree.IsLeaf() {
		return thisScore
	}
	leftScore := b.initialPrune(tree.LeftChild)
	rightScore := b.initialPrune(tree.RightChild)
	if thisScore < leftScore+rightScore+epsilon {
		// collapse this node
		tree.LeftChild = nil
		tree.RightChild = nil
	} else {
		thisScore = leftScore + rightScore
	}
	return thisScore
}

// getCandidates returns the candidate trees for pruning, and the
// corresponding critical values for the tuning parameter alpha.
func (b *Factory) getCandidates(tree *Tree, classes int) ([]*Tree, []float64) {
	// TODO(voss): the repeated calls to findWeakestLink compute many
	// link strengths repeatedly.  Is it worth the effort to cache the
	// results (and update along the path toward the root when
	// pruning)?
	b.initialPrune(tree)
	candidates := []*Tree{tree}
	alpha := []float64{0.0}
	for !tree.IsLeaf() {
		ctx := &pruneCtx{
			lowestPenalty: math.Inf(+1),
			pruneScore:    b.PruneScore,
		}
		ctx.findWeakestLink(tree, nil)
		tree = collapseSubtree(tree, ctx.bestPath)
		k := len(candidates)
		if ctx.lowestPenalty > alpha[k-1] {
			candidates = append(candidates, tree)
			alpha = append(alpha, ctx.lowestPenalty)
		} else {
			candidates[k-1] = tree
		}
	}
	alpha = append(alpha, math.Inf(+1))
	return candidates, alpha
}

type direction uint8

const (
	left direction = iota
	right
)

type pruneCtx struct {
	pruneScore impurity.Function

	lowestPenalty float64
	bestPath      []direction
}

// The findWeakestLink method updates the fields in `ctx` to indicate
// the "weakest link" in the tree.
//
// When called from the outside, the `path` argument must be nil (it
// is used internally in recursive calls).
//
// The method returns the total impurity and number of leaves of the
// subtree `t`.
func (ctx *pruneCtx) findWeakestLink(t *Tree, path []direction) (float64, int) {
	collapsedImpurity := ctx.pruneScore(t.Hist)
	if t.IsLeaf() {
		return collapsedImpurity, 1
	}

	leftImpurity, leftLeaves :=
		ctx.findWeakestLink(t.LeftChild, append(path, left))
	rightImpurity, rightLeaves :=
		ctx.findWeakestLink(t.RightChild, append(path, right))
	fullImpurity := leftImpurity + rightImpurity
	allLeaves := leftLeaves + rightLeaves

	penalty := (collapsedImpurity - fullImpurity) / float64(allLeaves-1)
	if penalty <= ctx.lowestPenalty {
		ctx.lowestPenalty = penalty
		ctx.bestPath = make([]direction, len(path))
		copy(ctx.bestPath, path)
	}

	return fullImpurity, allLeaves
}

// The collapseSubtree function returns a new tree with the subtree
// rooted at a given node collapsed into a single leaf node.  Internal
// nodes are copied as needed to ensure that the original tree is not
// modified by this procedure.
func collapseSubtree(tree *Tree, path []direction) *Tree {
	n := len(path)
	spine := make([]*Tree, n)
	t := tree
	for i, dir := range path {
		spine[i] = t
		if dir == left {
			t = t.LeftChild
		} else {
			t = t.RightChild
		}
	}
	res := &Tree{
		Hist: t.Hist,
	}
	for i := n - 1; i >= 0; i-- {
		if path[i] == left {
			res = &Tree{
				Hist:       spine[i].Hist,
				LeftChild:  res,
				RightChild: spine[i].RightChild,
				Column:     spine[i].Column,
				Limit:      spine[i].Limit,
			}
		} else {
			res = &Tree{
				Hist:       spine[i].Hist,
				LeftChild:  spine[i].LeftChild,
				RightChild: res,
				Column:     spine[i].Column,
				Limit:      spine[i].Limit,
			}
		}
	}
	return res
}

func selectCandidate(alpha []float64, a float64) int {
	for i := 0; i < len(alpha)-1; i++ {
		if alpha[i+1] >= a {
			return i
		}
	}
	panic("alpha not terminated by +inf")
}
