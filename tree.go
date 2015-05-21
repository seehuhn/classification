package classification

import (
	"fmt"
	"github.com/gonum/matrix/mat64"
	"sort"
	"strings"
)

type Tree struct {
	// fields used for internal nodes
	leftChild  *Tree
	rightChild *Tree
	col        int
	limit      float64

	// fields used for leaf nodes
	counts []int
}

func (t *Tree) doFormat(indent int) []string {
	pfx := strings.Repeat("  ", indent)
	res := []string{}
	if t.leftChild != nil {
		res = append(res, pfx+fmt.Sprintf("if x[%d] <= %g:", t.col, t.limit))
		res = append(res, t.leftChild.doFormat(indent+1)...)
		res = append(res, pfx+"else:")
		res = append(res, t.rightChild.doFormat(indent+1)...)
	} else {
		res = append(res, pfx+fmt.Sprintf("%v", t.counts))
	}
	return res
}

func (t *Tree) String() string {
	nodes := 0
	maxDepth := 0
	t.foreachLeaf(func(_ *Tree, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		nodes++
	}, 0)
	return fmt.Sprintf("<classification tree, %d leaves, max depth %d>",
		nodes, maxDepth)
}

func (t *Tree) Format() string {
	return strings.Join(t.doFormat(0), "\n")
}

func (t *Tree) Lookup(x []float64) []float64 {
	for {
		if t.leftChild == nil {
			return probabilities(t.counts)
		}
		if x[t.col] <= t.limit {
			t = t.leftChild
		} else {
			t = t.rightChild
		}
	}
}

func (t *Tree) costComplexityScore(score ScoreFunction, alpha float64) float64 {
	var res float64
	leaves := 0
	t.foreachLeaf(func(t *Tree, _ int) {
		leaves++
		res += float64(sum(t.counts)) * score(t.counts)
	}, 0)
	res += alpha * float64(leaves)
	return res
}

func (t *Tree) walkPostOrder(fn func(*Tree, int), depth int) {
	if t.leftChild != nil {
		t.leftChild.walkPostOrder(fn, depth+1)
		t.rightChild.walkPostOrder(fn, depth+1)
	}
	fn(t, depth)
}

func (t *Tree) foreachLeaf(fn func(*Tree, int), depth int) {
	if t.leftChild != nil {
		t.leftChild.foreachLeaf(fn, depth+1)
		t.rightChild.foreachLeaf(fn, depth+1)
	} else {
		fn(t, depth)
	}
}

type pruneCtx struct {
	pruneScore ScoreFunction

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
				col:        spine[i].col,
				limit:      spine[i].limit,
			}
		} else {
			res = &Tree{
				leftChild:  spine[i].leftChild,
				rightChild: res,
				col:        spine[i].col,
				limit:      spine[i].limit,
			}
		}
	}
	return res
}

type ScoreFunction func([]int) float64

type StopFunction func([]int) bool

type TreeBuilder struct {
	StopGrowth StopFunction
	SplitScore ScoreFunction
	PruneScore ScoreFunction
}

func (b *TreeBuilder) NewTree(x *mat64.Dense, k Classes, response []int, alpha float64) *Tree {
	// build the initial tree
	rows := make([]int, len(response))
	for i := range rows {
		rows[i] = i
	}
	xb := &xBuilder{*b, x, k, response}
	tree := xb.build(rows)

	// get all candidates for pruning the tree
	candidates := []*Tree{}
	for {
		candidates = append(candidates, tree)

		if tree.leftChild == nil {
			break
		}

		childCount := make([]int, int(k))
		ctx := &pruneCtx{
			pruneScore: b.PruneScore,
		}
		ctx.search(tree, nil, childCount)
		tree = ctx.prunedTree(tree)
	}

	// get the "optimal" pruned tree
	bestTree := candidates[0]
	bestScore := bestTree.costComplexityScore(b.PruneScore, alpha)
	for _, tree := range candidates[1:] {
		score := tree.costComplexityScore(b.PruneScore, alpha)
		if score < bestScore {
			bestTree = tree
			bestScore = score
		}
	}

	return bestTree
}

type xBuilder struct {
	TreeBuilder
	x        *mat64.Dense
	k        Classes
	response []int
}

func (b *xBuilder) build(rows []int) *Tree {
	y := make([]int, len(rows))
	for i, row := range rows {
		y[i] = b.response[row]
	}
	if b.StopGrowth(y) {
		return &Tree{
			counts: b.k.Frequencies(y),
		}
	}

	first := true
	var bestCol int
	var bestSplit int
	var bestLimit float64
	var bestScore float64

	_, nCol := b.x.Dims()
	for col := 0; col < nCol; col++ {
		sort.Sort(colSort{b.x, rows, col})
		sortedResp := make([]int, len(rows))
		for i, row := range rows {
			sortedResp[i] = b.response[row]
		}

		for i := 1; i < len(rows); i++ {
			limit := (b.x.At(rows[i-1], col) + b.x.At(rows[i], col)) / 2
			leftScore := b.SplitScore(b.k.Frequencies(sortedResp[:i]))
			rightScore := b.SplitScore(b.k.Frequencies(sortedResp[i:]))
			p := float64(i) / float64(len(rows))
			score := leftScore*p + rightScore*(1-p)
			if first || score < bestScore {
				bestCol = col
				bestSplit = i
				bestLimit = limit
				bestScore = score
				first = false
			}
		}
	}

	sort.Sort(colSort{b.x, rows, bestCol})
	return &Tree{
		leftChild:  b.build(rows[:bestSplit]),
		rightChild: b.build(rows[bestSplit:]),
		col:        bestCol,
		limit:      bestLimit,
	}
}

type colSort struct {
	x    mat64.Matrix
	rows []int
	col  int
}

func (c colSort) Len() int { return len(c.rows) }
func (c colSort) Less(i, j int) bool {
	return c.x.At(c.rows[i], c.col) < c.x.At(c.rows[j], c.col)
}
func (c colSort) Swap(i, j int) {
	c.rows[i], c.rows[j] = c.rows[j], c.rows[i]
}
