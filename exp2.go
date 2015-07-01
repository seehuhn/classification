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
	m, err := matrix.ReadAsText(trainFile, matrix.PlainFormat)
	if err != nil {
		fmt.Printf("cannot read %s: %s\n", trainFile, err.Error())
	}
	fmt.Println(m)
}
