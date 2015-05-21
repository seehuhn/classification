// +build ignore

package main

import (
	"fmt"
	"github.com/gonum/matrix/mat64"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/mt19937"
	"math/rand"
)

func main() {
	rng := rand.New(mt19937.New())
	rng.Seed(12)

	means := []float64{-1.0, 1.0}

	k := classification.Classes(2)

	n := 10000
	p := 2
	raw := make([]float64, n*p)
	response := make([]int, n)
	for i := 0; i < n; i++ {
		response[i] = rng.Intn(int(k))
		for j := 0; j < p; j++ {
			raw[i*p+j] = rng.NormFloat64() + means[response[i]]
		}
	}
	x := mat64.NewDense(n, p, raw)

	b := &classification.TreeBuilder{
		StopGrowth: func(y []int) bool {
			if len(y) <= 5 {
				return true
			}
			for i := 1; i < len(y); i++ {
				if y[i] != y[i-1] {
					return false
				}
			}
			return true
		},
		SplitScore: classification.Entropy,
		PruneScore: classification.MisclassificationError,
	}

	tree := b.NewTree(x, k, response, 3.0)
	fmt.Println(tree.Format())
}
