package classification

import (
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/util"
	"math"
)

const epsilon = 1e-6

// The initialPrune method modifies the given tree by recursively
// collapsing all leaves where the impurity `b.PruneScore` is not
// increased in the process.  The return value is the total impurity
// value of the pruned tree.
func (b *TreeBuilder) initialPrune(tree *Tree) float64 {
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

func (b *TreeBuilder) getCandidates(tree *Tree, classes int) []*Tree {
	candidates := []*Tree{tree}
	// TODO(voss): the repeated calls to findWeakestLink compute the
	// many link strengths repeatedly.  Is it worth the effort to
	// cache the (and update on the path toward the root when
	// pruining) results?
	for !tree.IsLeaf() {
		ctx := &pruneCtx{
			lowestPenalty: math.Inf(+1),
			pruneScore:    b.PruneScore,
		}
		ctx.findWeakestLink(tree, nil)
		tree = collapseSubtree(tree, ctx.bestPath)
		candidates = append(candidates, tree)
	}
	return candidates
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

func (t *Tree) costComplexityScore(score impurity.Function) (loss float64, leaves int) {
	t.ForeachLeaf(func(hist util.Histogram, _ int) {
		leaves++
		loss += score(hist)
	})
	return
}

func (b *TreeBuilder) filterCandidates(candidates []*Tree) ([]float64, []*Tree) {
	m := len(candidates)
	loss := make([]float64, m)
	complexity := make([]int, m)
	for i, tree := range candidates {
		loss[i], complexity[i] = tree.costComplexityScore(b.PruneScore)
	}

	breaks := []float64{}
	trees := []*Tree{}

	// Loss is increasing, complexity is decreasing.  We are looking
	// for the i which minimises loss[i] + alpha*complexity[i].
	// We find:
	//
	//   i is favoured over j
	//     <=> loss[i] + alpha*complexity[i] < loss[j] + alpha*complexity[j]
	//     <=> alpha*(complexity[i]-complexity[j]) < loss[j]-loss[i]
	//     <=> alpha < (loss[j]-loss[i]) / (complexity[i]-complexity[j])
	i := 0
	for i < m {
		bestJ := m
		bestStop := math.Inf(+1)
		bestVal := math.Inf(+1)
		for j := i + 1; j < m && loss[j] < bestVal; j++ {
			stop := (loss[j] - loss[i]) / float64(complexity[i]-complexity[j])
			if stop <= bestStop {
				bestJ = j
				bestStop = stop
				bestVal = loss[i] + stop*float64(complexity[i])
			}
		}
		if bestJ < m {
			breaks = append(breaks, bestStop)
		}
		trees = append(trees, candidates[i])
		i = bestJ
	}
	return breaks, trees
}

func (b *TreeBuilder) selectCandidate(candidates []*Tree, alpha float64) *Tree {
	bestPenalty := math.Inf(+1)
	var bestTree *Tree
	for _, tree := range candidates {
		loss, complexity := tree.costComplexityScore(b.PruneScore)
		penalty := loss + alpha*float64(complexity)
		if penalty <= bestPenalty {
			bestPenalty = penalty
			bestTree = tree
		}
	}
	return bestTree
}
