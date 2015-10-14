// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/matrix"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func zipDataColumns(col int) matrix.ColumnType {
	switch col {
	case 0:
		return matrix.IntColumn
	default:
		return matrix.Float64Column
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

	trainFile := "zip.train.gz"
	XTrain, YTrain, err := matrix.Plain.Read(trainFile, zipDataColumns)
	if err != nil {
		fmt.Printf("cannot read %s: %s\n", trainFile, err.Error())
	}

	testFile := "zip.test.gz"
	XTest, YTest, err := matrix.Plain.Read(testFile, zipDataColumns)
	if err != nil {
		fmt.Printf("cannot read %s: %s\n", testFile, err.Error())
	}

	b := &classification.TreeBuilder{
		XValLoss:   loss.ZeroOne,
		SplitScore: impurity.Entropy,
		PruneScore: impurity.Entropy,
		K:          2,
	}
	tree, est := b.TreeFromTrainingsData(10, XTrain, YTrain.Column(0))
	fmt.Println(tree)
	fmt.Println("estimated average loss from cross validation", est)

	n, _ := XTest.Shape()
	sum := 0.0
	wrong := 0
	for i := 0; i < n; i++ {
		correct := YTest.At(i, 0)
		row := XTest.Row(i)
		hist := tree.EstimateClassProbabilities(row)
		sum += b.XValLoss(correct, hist)
		if tree.GuessClass(row) != correct {
			wrong++
		}
	}
	fmt.Println("average loss from test set", sum/float64(n))
	fmt.Printf("misclassification rate for test set: %d / %d = %g\n",
		wrong, n, float64(wrong)/float64(n))
}
