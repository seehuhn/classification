// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/matrix"
)

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
		SplitScore: impurity.Entropy,
	}
	tree, est := b.NewTree(XTrain, 10, YTrain.Column(0))
	fmt.Println(tree, est)
}
