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
	XTrain, YTrain, err := matrix.ReadAsText(trainFile,
		zipDataColumns, matrix.PlainFormat)
	if err != nil {
		fmt.Printf("cannot read %s: %s\n", trainFile, err.Error())
	}

	testFile := "zip.test.gz"
	XTest, YTest, err := matrix.ReadAsText(testFile,
		zipDataColumns, matrix.PlainFormat)
	if err != nil {
		fmt.Printf("cannot read %s: %s\n", testFile, err.Error())
	}

	b := &classification.TreeBuilder{
		XValLoss:   loss.Other,
		SplitScore: impurity.Entropy,
		K:          2,
	}
	tree, est := b.NewTree(XTrain, 10, YTrain.Column(0))
	fmt.Println(tree, est)

	n, _ := XTest.Shape()
	sum := 0.0
	for i := 0; i < n; i++ {
		correct := YTest.At(i, 0)
		row := XTest.Row(i)
		hist := tree.Lookup(row)
		sum += b.XValLoss(correct, hist)
	}
	fmt.Println("test set:", sum/float64(n))
}
