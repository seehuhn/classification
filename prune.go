package classification

import (
	"fmt"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/util"
	"math"
)

func (b *TreeBuilder) prunedTrees(tree *Tree, classes int) []*Tree {
	candidates := []*Tree{tree}
	fmt.Println(tree.Format())
	for tree.LeftChild != nil {
		ctx := &pruneCtx{
			lowestPenalty: math.Inf(1),
			pruneScore:    b.PruneScore,
		}
		ctx.findWeakestLink(tree, nil)
		tree = collapseSubtree(tree, ctx.bestPath)
		candidates = append(candidates, tree)
	}
	fmt.Println("----------------------")
	return candidates
}

type pruneCtx struct {
	pruneScore impurity.Function

	lowestPenalty float64
	bestPath      []direction
}

// The findWeakestLink method updates the fields in `ctx` to indicate
// the link in the tree which contributes least to the overall `pruneScore`.
//
// When called from the outside, the `path` argument must be nil (it
// is used internally in recursive calls).
//
// The method returns the total score of the full subtree.
func (ctx *pruneCtx) findWeakestLink(t *Tree, path []direction) float64 {
	collapsedScore := ctx.pruneScore(t.Hist)
	if t.LeftChild == nil {
		return collapsedScore
	}

	leftFullScore := ctx.findWeakestLink(t.LeftChild, append(path, left))
	rightFullScore := ctx.findWeakestLink(t.RightChild, append(path, right))
	fullScore := leftFullScore + rightFullScore

	penalty := collapsedScore - fullScore
	fmt.Println(path, fullScore, collapsedScore, penalty)
	if penalty < ctx.lowestPenalty {
		ctx.lowestPenalty = penalty
		ctx.bestPath = make([]direction, len(path))
		copy(ctx.bestPath, path)
	}

	return fullScore
}

// The collapseSubtree function returns a new tree with the subtree rooted at
// a given node collapsed into a single leaf node.  Internal nodes are
// copied as needed to ensure that the original tree is not modified
// by this procedure.
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

type direction uint8

const (
	left direction = iota
	right
)

func (t *Tree) costComplexityScore(score impurity.Function) (loss float64, leaves int) {
	t.ForeachLeaf(func(hist util.Histogram, _ int) {
		leaves++
		loss += score(hist)
	})
	return
}

func (b *TreeBuilder) tryTrees(candidates []*Tree, alpha []float64) []*Tree {
	// get expected loss and complexity
	// TODO(voss): ctx.findWeakestLink, above, already computes the expected
	// loss; reuse these values rather than computing them again?
	m := len(candidates)
	loss := make([]float64, m)
	complexity := make([]int, m)
	for i, tree := range candidates {
		loss[i], complexity[i] = tree.costComplexityScore(b.PruneScore)
	}

	// Loss is increasing, complexity is decreasing.  We are looking
	// for the i which minimises loss[i] + alpha*complexity[i].
	// We find:
	//
	//   i is favoured over i+1
	//     <=> loss[i] + alpha*complexity[i] < loss[i+1] + alpha*complexity[i+1]
	//     <=> alpha*(complexity[i]-complexity[i+1]) < loss[i+1]-loss[i]
	//     <=> alpha < (loss[i+1]-loss[i]) / (complexity[i]-complexity[i+1])
	if alpha[0] < 0 {
		var alphaMin float64
		i := 0
		for alphaMin < 1e-6 {
			alphaMin = 0.8 * (loss[i+1] - loss[i]) /
				float64(complexity[i]-complexity[i+1])
			i++
		}
		alphaMax := 1.2 * (loss[m-1] - loss[m-2]) /
			float64(complexity[m-2]-complexity[m-1])
		alphaSteps := len(alpha)
		step := math.Pow(alphaMax/alphaMin, 1/float64(alphaSteps-1))
		a := alphaMin
		for i := range alpha {
			alpha[i] = a
			a *= step
		}
	}

	// get the "optimal" pruned tree
	bestScore := make([]float64, len(alpha))
	for i := range bestScore {
		bestScore[i] = math.Inf(1)
	}
	bestTree := make([]*Tree, len(alpha))
	bestIdx := make([]int, len(alpha))
	for i, a := range alpha {
		for j, tree := range candidates {
			score := loss[j] + a*float64(complexity[j])
			if score < bestScore[i] {
				bestScore[i] = score
				bestTree[i] = tree
				bestIdx[i] = j
			}
		}
	}

	return bestTree
}
