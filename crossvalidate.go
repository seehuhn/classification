package classification

// NewTree constructs a new classification tree.
//
// K-fold crossvalidation is used to find the optimal pruning
// parameter.
func (b *TreeBuilder) NewTree(x *Matrix, classes int, response []int) (*Tree, float64) {

	n := len(response)
	learnSize := n * (b.K - 1) / b.K

	alphaSteps := 50
	alpha := make([]float64, alphaSteps)
	alpha[0] = -1 // ask tryTrees to determine the range of alpha

	mean := make([]float64, len(alpha))

	for k := 0; k < b.K; k++ {
		learnRows := make([]int, 0, learnSize+1)
		testRows := make([]int, 0, n-learnSize)
		for i := range response {
			if i%b.K == k {
				testRows = append(testRows, i)
			} else {
				learnRows = append(learnRows, i)
			}
		}

		trees := b.tryTrees(x, classes, response, learnRows, alpha)

		cache := make(map[*Tree][2]float64)
		for l, tree := range trees {
			var cumLoss [2]float64
			if loss, ok := cache[tree]; ok {
				cumLoss = loss
			} else {
				for _, row := range testRows {
					prob := tree.Lookup(x.Row(row))
					val := b.XValLoss(response[row], prob)
					cumLoss[0] += val
					cumLoss[1] += val * val
				}
				cache[tree] = cumLoss
			}
			mean[l] += cumLoss[0]
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

	rows := make([]int, len(response))
	for i := range rows {
		rows[i] = i
	}
	tree := b.tryTrees(x, classes, response, rows, []float64{bestAlpha})[0]

	return tree, bestExpectedLoss
}
