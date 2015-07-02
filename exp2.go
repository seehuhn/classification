// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/seehuhn/classification/matrix"
)

func main() {
	flag.Parse()

	trainFile := "zip.train.gz"
	X, Y, err := matrix.ReadAsText(trainFile,
		func(col int) matrix.ColumnType {
			switch col {
			case 0:
				return matrix.IntColumn
			default:
				return matrix.Float64Column
			}
		},
		matrix.PlainFormat)
	if err != nil {
		fmt.Printf("cannot read %s: %s\n", trainFile, err.Error())
	}
	fmt.Println(X)
	fmt.Println(Y)
}
