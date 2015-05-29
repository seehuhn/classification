package classification

import (
	"github.com/seehuhn/classification/impurity"
	"math"
)

func (b *TreeBuilder) prunedTrees(tree *Tree, classes int) []*Tree {
	candidates := []*Tree{}
	for {
		candidates = append(candidates, tree)

		if tree.leftChild == nil {
			break
		}

		ctx := &pruneCtx{
			lowestPenalty: math.Inf(1),
			pruneScore:    b.PruneScore,
		}
		ctx.search(tree, nil)
		tree = ctx.prunedTree(tree)
	}
	return candidates
}

func (t *Tree) costComplexityScore(score impurity.Function) (loss float64, leaves int) {
	t.foreachLeaf(func(t *Tree, _ int) {
		leaves++
		loss += score(t.counts)
	}, 0)
	return
}

// func (b *TreeBuilder) tryTrees(x *Matrix, classes int, response []int, rows []int, alpha []float64) []*Tree {
func (b *TreeBuilder) tryTrees(candidates []*Tree, alpha []float64) []*Tree {

	// get expected loss and complexity
	// TODO(voss): ctx.search, above, already computes the expected
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

type direction uint8

const (
	left direction = iota
	right
)

type pruneCtx struct {
	pruneScore impurity.Function

	lowestPenalty float64
	bestPath      []direction
	bestCounts    []int
}

// returns the total score of the full subtree
func (ctx *pruneCtx) search(t *Tree, path []direction) float64 {
	collapsedScore := ctx.pruneScore(t.counts)
	if t.leftChild == nil {
		return collapsedScore
	}

	leftFullScore := ctx.search(t.leftChild, append(path, left))
	rightFullScore := ctx.search(t.rightChild, append(path, right))
	fullScore := leftFullScore + rightFullScore

	penalty := collapsedScore - fullScore
	if penalty < ctx.lowestPenalty {
		ctx.lowestPenalty = penalty
		ctx.bestPath = make([]direction, len(path))
		copy(ctx.bestPath, path)
		ctx.bestCounts = t.counts
	}

	return fullScore
}

func (ctx *pruneCtx) prunedTree(tree *Tree) *Tree {
	n := len(ctx.bestPath)
	spine := make([]*Tree, n)
	t := tree
	for i, dir := range ctx.bestPath {
		spine[i] = t
		if dir == left {
			t = t.leftChild
		} else {
			t = t.rightChild
		}
	}
	res := &Tree{
		counts: ctx.bestCounts,
	}
	for i := n - 1; i >= 0; i-- {
		if ctx.bestPath[i] == left {
			res = &Tree{
				leftChild:  res,
				rightChild: spine[i].rightChild,
				column:     spine[i].column,
				limit:      spine[i].limit,
				counts:     spine[i].counts,
			}
		} else {
			res = &Tree{
				leftChild:  spine[i].leftChild,
				rightChild: res,
				column:     spine[i].column,
				limit:      spine[i].limit,
				counts:     spine[i].counts,
			}
		}
	}
	return res
}
