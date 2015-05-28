package classification

import (
	"github.com/seehuhn/classification/impurity"
	"math"
)

func (t *Tree) costComplexityScore(score impurity.Function) (loss float64, leaves int) {
	t.foreachLeaf(func(t *Tree, _ int) {
		leaves++
		loss += float64(sum(t.counts)) * score(t.counts)
	}, 0)
	return
}

func (b *TreeBuilder) tryTrees(x *Matrix, classes int, response []int,
	rows []int, alpha []float64) []*Tree {

	// build the initial tree
	xb := &xBuilder{*b, x, classes, response}
	tree := xb.build(rows)

	// get all candidates for pruning the tree
	candidates := []*Tree{}
	for {
		candidates = append(candidates, tree)

		if tree.leftChild == nil {
			break
		}

		ctx := &pruneCtx{
			pruneScore: b.PruneScore,
		}
		childCount := make([]int, classes)
		ctx.search(tree, nil, childCount)
		tree = ctx.prunedTree(tree)
	}

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

type pruneCtx struct {
	pruneScore impurity.Function

	valid         bool
	lowestPenalty float64
	bestPath      []bool
	bestCounts    []int
}

func (ctx *pruneCtx) search(t *Tree, path []bool, childCount []int) float64 {
	if t.leftChild == nil {
		for i, ni := range t.counts {
			childCount[i] += ni
		}
		return float64(sum(t.counts)) * ctx.pruneScore(t.counts)
	}

	var res float64
	childSubCount := make([]int, len(childCount))
	res += ctx.search(t.leftChild, append(path, true), childSubCount)
	res += ctx.search(t.rightChild, append(path, false), childSubCount)
	for i, ni := range childSubCount {
		childCount[i] += ni
	}

	local := float64(sum(childSubCount)) * ctx.pruneScore(childSubCount)
	penalty := local - res
	if !ctx.valid || penalty < ctx.lowestPenalty {
		ctx.valid = true
		ctx.lowestPenalty = penalty
		ctx.bestPath = make([]bool, len(path))
		copy(ctx.bestPath, path)
		ctx.bestCounts = childSubCount
	}
	return res
}

func (ctx *pruneCtx) prunedTree(tree *Tree) *Tree {
	n := len(ctx.bestPath)
	spine := make([]*Tree, n)
	t := tree
	for i, left := range ctx.bestPath {
		spine[i] = t
		if left {
			t = t.leftChild
		} else {
			t = t.rightChild
		}
	}
	res := &Tree{
		counts: ctx.bestCounts,
	}
	for i := n - 1; i >= 0; i-- {
		if ctx.bestPath[i] {
			res = &Tree{
				leftChild:  res,
				rightChild: spine[i].rightChild,
				column:     spine[i].column,
				limit:      spine[i].limit,
			}
		} else {
			res = &Tree{
				leftChild:  spine[i].leftChild,
				rightChild: res,
				column:     spine[i].column,
				limit:      spine[i].limit,
			}
		}
	}
	return res
}
