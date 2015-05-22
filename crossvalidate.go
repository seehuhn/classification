package classification

import (
	"fmt"
	"github.com/gonum/matrix/mat64"
	"math"
)

// NewXVTree constructs a new classification tree using K-fold
// crossvalidation to find the optimal pruning parameter alpha.
func (b *TreeBuilder) NewXVTree(x *mat64.Dense, classes Classes, response []int, K int) *Tree {
	n := len(response)
	learnSize := n * (K - 1) / K

	alphaSteps := 50
	alpha := make([]float64, alphaSteps)
	alpha[0] = -1 // ask NewTree to get the range of alpha

	mean := make([]float64, len(alpha))
	squares := make([]float64, len(alpha))

	for k := 0; k < K; k++ {
		learnRows := make([]int, 0, learnSize+1)
		testRows := make([]int, 0, n-learnSize)
		for i := range response {
			if i%K == k {
				testRows = append(testRows, i)
			} else {
				learnRows = append(learnRows, i)
			}
		}

		trees := b.NewTrees(x, classes, response, learnRows, alpha)

		cache := make(map[*Tree][2]float64)
		for l, tree := range trees {
			var cumLoss [2]float64
			if loss, ok := cache[tree]; ok {
				cumLoss = loss
			} else {
				for _, row := range testRows {
					prob := tree.Lookup(x.Row(nil, row))
					val := b.XValLoss(response[row], prob)
					cumLoss[0] += val
					cumLoss[1] += val * val
				}
				cache[tree] = cumLoss
			}
			mean[l] += cumLoss[0]
			squares[l] += cumLoss[1]
		}
	}
	for l := range alpha {
		mean[l] /= float64(n)
		squares[l] /= float64(n)
	}

	for l, a := range alpha {
		m := mean[l]
		sd := math.Sqrt(squares[l] - m*m)
		fmt.Printf("%f %f %f\n", a, m, sd/math.Sqrt(float64(n)))
	}

	return nil
}
