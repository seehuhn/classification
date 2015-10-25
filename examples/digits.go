// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/tree"
	"log"
	"math"
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

	digits := data.Digits

	b := &tree.Factory{
		XValLoss:   loss.ZeroOne,
		SplitScore: impurity.Gini,
		PruneScore: impurity.MisclassificationError,
		K:          5,
	}
	trainingData, err := digits.TrainingData()
	if err != nil {
		log.Fatal(err)
	}
	tree := b.FromData(trainingData)
	fmt.Println(tree)

	testData, err := digits.TestData()
	if err != nil {
		log.Fatal(err)
	}
	cumLoss := 0.0
	cumLoss2 := 0.0
	nTest := len(testData.Y)
	for i := 0; i < nTest; i++ {
		row := testData.X.Row(i)
		prob := tree.EstimateClassProbabilities(row)
		l := b.XValLoss(testData.Y[i], prob)
		cumLoss += l
		cumLoss2 += l * l
	}
	nn := float64(nTest)
	cumLoss /= nn
	cumLoss2 /= nn
	stdErr := math.Sqrt((cumLoss2 - cumLoss*cumLoss) / nn)

	fmt.Println("average loss from test set:", cumLoss, "+-", stdErr)
}
