// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/mt19937"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rng := rand.New(mt19937.New())
	rng.Seed(12)

	classes := 2
	means := []float64{-1, 1}

	n := 20000
	p := 2
	raw := make([]float64, n*p)
	response := make([]int, n)
	for i := 0; i < n; i++ {
		response[i] = rng.Intn(classes)
		for j := 0; j < p; j++ {
			raw[i*p+j] = rng.NormFloat64() + means[response[i]]
		}
	}
	x := classification.NewMatrix(n, p, raw)

	b := &classification.TreeBuilder{
		XValLoss: loss.Other,
		K:        5,

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
		SplitScore: impurity.Entropy,
		PruneScore: impurity.MisclassificationError,
	}

	_, estLoss := b.NewTree(x, classes, response)

	fmt.Println(n, estLoss)
}
