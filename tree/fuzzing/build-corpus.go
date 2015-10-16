package main

import (
	"fmt"
	"github.com/seehuhn/classification/tree"
	"math"
	"os"
	"path/filepath"
)

var fileIndex int

func writeTree(t *tree.Tree) {
	fileIndex++
	fname := filepath.Join("corpus", fmt.Sprintf("simple%02d.bin", fileIndex))
	fd, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	err = t.WriteBinary(fd)
	if err != nil {
		panic(err)
	}
	fd.Close()
}

func main() {
	for p := 1; p < 600; p *= 2 {
		t := &tree.Tree{
			Hist: make([]float64, p),
		}
		t.Hist[p/2] = float64(p)
		writeTree(t)
	}
	for p := 1; p < 600; p *= 2 {
		t := &tree.Tree{
			Hist: make([]float64, p),
			LeftChild: &tree.Tree{
				Hist: make([]float64, p),
			},
			RightChild: &tree.Tree{
				Hist: make([]float64, p),
			},
			Column: p / 2,
			Limit:  math.Sin(float64(p)) * float64(p),
		}
		t.LeftChild.Hist[0] = 100
		t.RightChild.Hist[0] = 200
		t.Hist[0] = 300
		writeTree(t)
	}
}
