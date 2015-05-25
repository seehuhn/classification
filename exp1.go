// files.go -
// Copyright (C) 2015  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var rng *rand.Rand

func sample() (x float64, y int) {
	x = rng.Float64()
	p := (1 + 2*x) / 4
	if rng.Float64() < p {
		y = 1
	}
	return
}

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

	rng = rand.New(mt19937.New())
	rng.Seed(1)

	n := 5000
	raw := make([]float64, n)
	response := make([]int, n)
	for i := 0; i < n; i++ {
		raw[i], response[i] = sample()
	}
	x := classification.NewMatrix(n, 1, raw)

	b := &classification.TreeBuilder{
		XValLoss: loss.Other,
		K:        5,

		StopGrowth: func(y []int) bool {
			return len(y) <= 5
		},
		SplitScore: impurity.Entropy,
		PruneScore: impurity.MisclassificationError,
	}

	tree, estLoss := b.NewTree(x, 2, response)

	N := 1000
	var lVal, lSquaredVal float64
	for j := 0; j < N; j++ {
		xj, yj := sample()
		pj := tree.Lookup([]float64{xj})
		l := b.XValLoss(yj, pj)
		lVal += l
		lSquaredVal += l * l
	}
	lVal /= float64(N)
	lSquaredVal /= float64(N)

	fmt.Println(tree.Format())

	fmt.Println(n, estLoss, lVal, math.Sqrt((lSquaredVal-lVal*lVal)/float64(N)))
}
