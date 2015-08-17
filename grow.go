package classification

import (
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/stop"
	"github.com/seehuhn/classification/util"
	"sort"
)

// TreeBuilder is a structure to store the parameters governing the
// growing and pruning of classification trees.  Any zero field values
// are interpreted as the corresponding values from the
// `DefaultTreeBuilder` structure.
type TreeBuilder struct {
	// StopGrowth decides for every node of the initial tree whether a
	// further split is considered.  The default is to stop splitting
	// nodes once the node only contains a single type of
	// observations.
	StopGrowth stop.Function

	// SplitScore is the impurity function used when growing the
	// initial tree.  Splits are (greedily) chosen to minimise the
	// total `SplitScore`.  The default is to use Gini impurity here.
	SplitScore impurity.Function

	// PruneScore is the cost function used for cost-complexity
	// pruning.  Pruning aims to reduce the size of the tree (in order
	// to avoid overfitting), while keeping the `PruneScore` small.
	// The default is to use the misclassification error rate for
	// pruning.
	PruneScore impurity.Function

	// XValLoss is used to guide the balance between the cost (as
	// measured by `PruneScore`) and complexity (i.e. the final size
	// of the tree).  Cost and complexity are balanced to minimise
	// XValLoss.  The default is to use zero-one loss.
	XValLoss loss.Function

	// The number of groups to use in cross-validation when estimating
	// the expected loss.  The default is to use 5 groups.
	K int
}

// DefaultTreeBuilder specifies the default parameters for
// constructing a tree; see the `TreeBuilder` documentation for the
// meaning of the individual fields.  The values given in
// `DefaultTreeBuilder` are used by the `TreeFromTrainingsData`
// function, and to replace zero values in a `TreeBuilder` structure
// when the `TreeBuilder.TreeFromTrainingsData` method is called.
var DefaultTreeBuilder = &TreeBuilder{
	StopGrowth: stop.IfHomogeneous,
	SplitScore: impurity.Gini,
	PruneScore: impurity.MisclassificationError,
	XValLoss:   loss.ZeroOne,
	K:          5,
}

func (b *TreeBuilder) setDefaults() {
	if b.XValLoss == nil {
		b.XValLoss = DefaultTreeBuilder.XValLoss
	}
	if b.K == 0 {
		b.K = DefaultTreeBuilder.K
	}
	if b.StopGrowth == nil {
		b.StopGrowth = DefaultTreeBuilder.StopGrowth
	}
	if b.SplitScore == nil {
		b.SplitScore = DefaultTreeBuilder.SplitScore
	}
	if b.PruneScore == nil {
		b.PruneScore = DefaultTreeBuilder.PruneScore
	}
}

func (b *TreeBuilder) fullTree(x *matrix.Float64, classes int, response []int) *Tree {
	b.setDefaults()
	rows := intRange(len(response))
	hist := util.GetHist(rows, classes, response)
	xb := &xBuilder{*b, x, classes, response}
	return xb.getFullTree(rows, hist)
}

// TreeFromTrainingsData constructs a new classification tree from
// trainings data.
//
// K-fold crossvalidation is used to find the optimal pruning
// parameter.
//
// The return values are the new tree and an estimate for the average
// value of the loss function (given by `b.XValLoss`).
func (b *TreeBuilder) TreeFromTrainingsData(classes int, x *matrix.Float64, response []int) (*Tree, float64) {
	b.setDefaults()

	n := len(response)

	alphaSteps := 50
	alpha := make([]float64, alphaSteps)
	alpha[0] = -1 // ask tryTrees to determine the range of alpha

	mean := make([]float64, len(alpha))

	for k := 0; k < b.K; k++ {
		learnRows, testRows := getXValSets(k, b.K, n)

		// build the initial tree
		learnHist := util.GetHist(learnRows, classes, response)
		xb := &xBuilder{*b, x, classes, response}
		tree := xb.getFullTree(learnRows, learnHist)
		b.initialPrune(tree)

		// get all candidates for pruning the tree
		candidates := b.prunedTrees(tree, classes)

		// get the optimally pruned tree for every alpha
		trees := b.tryTrees(candidates, alpha)

		cache := make(map[*Tree]float64)
		for l, tree := range trees {
			var cumLoss float64
			if loss, ok := cache[tree]; ok {
				cumLoss = loss
			} else {
				for _, row := range testRows {
					prob := tree.EstimateClassProbabilities(x.Row(row))
					val := b.XValLoss(response[row], prob)
					cumLoss += val
				}
				cache[tree] = cumLoss
			}
			mean[l] += cumLoss
		}
	}
	for l := range alpha {
		mean[l] /= float64(n)
	}

	var bestAlpha float64
	var bestExpectedLoss float64
	for l, a := range alpha {
		if l == 0 || mean[l] < bestExpectedLoss {
			bestAlpha = a
			bestExpectedLoss = mean[l]
		}
	}

	tree := b.fullTree(x, classes, response)
	b.initialPrune(tree)
	candidates := b.prunedTrees(tree, classes)
	tree = b.tryTrees(candidates, []float64{bestAlpha})[0]

	return tree, bestExpectedLoss
}

type xBuilder struct {
	TreeBuilder
	x        *matrix.Float64
	classes  int
	response []int
}

// plan to reduce the amount of sorting required
//
// 1. sort rows by col j: i0, i1, i2, ..., in only once
// 2. split rows: i0, ..., ik | i{k+1}, ..., in as before
// 3. for every other row:
//    - old sort order is: i0', i1', ..., in'
//    - after the split, the order stays the same, but elements are
//      sorted into two groups.

func (b *xBuilder) getFullTree(rows []int, hist util.Histogram) *Tree {
	if b.StopGrowth(hist) {
		return &Tree{
			Hist: hist,
		}
	}

	best := b.findBestSplit(rows, hist)

	return &Tree{
		Hist:       hist,
		LeftChild:  b.getFullTree(best.Left, best.LeftHist),
		RightChild: b.getFullTree(best.Right, best.RightHist),
		Column:     best.Col,
		Limit:      best.Limit,
	}
}

func (b *xBuilder) findBestSplit(rows []int, hist util.Histogram) *searchResult {
	// TODO(voss): notice if the split does not bring any improvement
	// and, in this case, do not split any further?
	best := &searchResult{}
	first := true
	_, p := b.x.Shape()
	for col := 0; col < p; col++ {
		rows = copyIntSlice(rows)
		sort.Sort(&colSort{b.x, rows, col})

		leftHist := make(util.Histogram, len(hist))
		var rightHist = copyIntSlice(hist)
		for i := 1; i < len(rows); i++ {
			yi := b.response[rows[i-1]]
			leftHist[yi]++
			rightHist[yi]--

			left := b.x.At(rows[i-1], col)
			right := b.x.At(rows[i], col)
			if !(left < right) {
				continue
			}
			limit := (left + right) / 2

			leftScore := b.SplitScore(leftHist)
			rightScore := b.SplitScore(rightHist)
			score := leftScore + rightScore

			if first || score < best.Score {
				best.Col = col
				best.Limit = limit
				best.Left = rows[:i]
				best.Right = rows[i:]
				best.LeftHist = copyIntSlice(leftHist)
				best.RightHist = copyIntSlice(rightHist)
				best.Score = score
				first = false
			}
		}
	}
	return best
}

type searchResult struct {
	Col                 int
	Limit               float64
	Left, Right         []int
	LeftHist, RightHist util.Histogram
	Score               float64
}

type colSort struct {
	x    *matrix.Float64
	rows []int
	col  int
}

func (c *colSort) Len() int { return len(c.rows) }
func (c *colSort) Less(i, j int) bool {
	return c.x.At(c.rows[i], c.col) < c.x.At(c.rows[j], c.col)
}
func (c *colSort) Swap(i, j int) {
	c.rows[i], c.rows[j] = c.rows[j], c.rows[i]
}
