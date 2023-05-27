// exp1.go -
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

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/seehuhn/mt19937"

	"seehuhn.de/go/classification/data"
	"seehuhn.de/go/classification/impurity"
	"seehuhn.de/go/classification/matrix"
	"seehuhn.de/go/classification/tree"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var rng *rand.Rand

func sample() (x float64, y int) {
	x = rng.Float64()
	p := 0.25 + 0.5*x
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
	rng.Seed(time.Now().UnixNano())

	n := 20
	raw := make([]float64, n)
	response := make([]int, n)
	for i := 0; i < n; i++ {
		raw[i], response[i] = sample()
	}
	d := &data.Data{
		NumClasses: 2,
		X:          matrix.NewFloat64(n, 1, 0, raw),
		Y:          response,
	}

	builder := &tree.Factory{
		PruneScore: impurity.Gini,
	}
	builder.FromData(d)
}
