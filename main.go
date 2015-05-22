// +build ignore

package main

import (
	"flag"
	"github.com/gonum/matrix/mat64"
	"github.com/seehuhn/classification"
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

	means := []float64{-1, 1}

	k := classification.Classes(2)

	n := 15000
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
	// classification.WriteMatrix("data", x)
	// classification.WriteVector("resp", response)

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
		XValLoss:   classification.OtherLoss,
	}

	b.NewXVTree(x, k, response, 10)

	tree := b.NewTree(x, k, response, 3.0)
	// fmt.Println(tree.Format())
	_ = tree
}
