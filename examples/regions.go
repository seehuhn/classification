// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/tree"
	"github.com/seehuhn/classification/util"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func fixHist(t *tree.Tree) {
	if !t.IsLeaf() {
		fixHist(t.LeftChild)
		fixHist(t.RightChild)
		t.Hist = make([]float64, 2)
		t.Hist[0] = t.LeftChild.Hist[0] + t.RightChild.Hist[0]
		t.Hist[1] = t.LeftChild.Hist[1] + t.RightChild.Hist[1]
	}
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// 1 2 3 4
	// 5 6 7 8
	// 9 0 1 2
	// 3 4 5 10
	model := &tree.Tree{
		Column: 0,
		Limit:  0.5,
		LeftChild: &tree.Tree{
			Column: 0,
			Limit:  0.25,
			LeftChild: &tree.Tree{
				// column 0
				Column: 1,
				Limit:  0.5,
				LeftChild: &tree.Tree{
					Column: 1,
					Limit:  0.25,
					LeftChild: &tree.Tree{
						// (0, 0)
						Hist: []float64{3, 7},
					},
					RightChild: &tree.Tree{
						// (0, 1)
						Hist: []float64{9, 1},
					},
				},
				RightChild: &tree.Tree{
					Column: 1,
					Limit:  0.75,
					LeftChild: &tree.Tree{
						// (0, 2)
						Hist: []float64{5, 5},
					},
					RightChild: &tree.Tree{
						// (0, 3)
						Hist: []float64{1, 9},
					},
				},
			},
			RightChild: &tree.Tree{
				// column 1
				Column: 1,
				Limit:  0.5,
				LeftChild: &tree.Tree{
					Column: 1,
					Limit:  0.25,
					LeftChild: &tree.Tree{
						// (1, 0)
						Hist: []float64{4, 6},
					},
					RightChild: &tree.Tree{
						// (1, 1)
						Hist: []float64{0, 10},
					},
				},
				RightChild: &tree.Tree{
					Column: 1,
					Limit:  0.75,
					LeftChild: &tree.Tree{
						// (1, 2)
						Hist: []float64{6, 4},
					},
					RightChild: &tree.Tree{
						// (1, 3)
						Hist: []float64{2, 8},
					},
				},
			},
		},
		RightChild: &tree.Tree{
			Column: 0,
			Limit:  0.75,
			LeftChild: &tree.Tree{
				// column 2
				Column: 1,
				Limit:  0.5,
				LeftChild: &tree.Tree{
					Column: 1,
					Limit:  0.25,
					LeftChild: &tree.Tree{
						// (2, 0)
						Hist: []float64{5, 5},
					},
					RightChild: &tree.Tree{
						// (2, 1)
						Hist: []float64{1, 9},
					},
				},
				RightChild: &tree.Tree{
					Column: 1,
					Limit:  0.75,
					LeftChild: &tree.Tree{
						// (2, 2)
						Hist: []float64{7, 3},
					},
					RightChild: &tree.Tree{
						// (2, 3)
						Hist: []float64{3, 7},
					},
				},
			},
			RightChild: &tree.Tree{
				// column 3
				Column: 1,
				Limit:  0.5,
				LeftChild: &tree.Tree{
					Column: 1,
					Limit:  0.25,
					LeftChild: &tree.Tree{
						// (3, 0)
						Hist: []float64{10, 0},
					},
					RightChild: &tree.Tree{
						// (3, 1)
						Hist: []float64{2, 8},
					},
				},
				RightChild: &tree.Tree{
					Column: 1,
					Limit:  0.75,
					LeftChild: &tree.Tree{
						// (3, 2)
						Hist: []float64{8, 2},
					},
					RightChild: &tree.Tree{
						// (3, 3)
						Hist: []float64{4, 6},
					},
				},
			},
		},
	}
	fixHist(model)

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	n := 100000
	p := 2
	raw := make([]float64, n*p)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		x0 := rng.Float64()
		x1 := rng.Float64()
		raw[2*i] = x0
		raw[2*i+1] = x1
		prob := model.EstimateClassProbabilities([]float64{x0, x1})
		U := rng.Float64()
		if U < prob[0] {
			y[i] = 0
		} else {
			y[i] = 1
		}
	}
	x := matrix.NewFloat64(n, p, 0, raw)

	tree, estLoss := tree.NewFromTrainingsData(2, x, y, nil)
	fmt.Println(estLoss)
	tree.ForeachLeafRegion(func(a, b []float64, hist util.Histogram, depth int) {
		fmt.Println(a[0], b[0], a[1], b[1], hist.Probabilities())
	})
}
