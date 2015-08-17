package main

import (
	"fmt"
	"github.com/seehuhn/classification"
	"math"
	"os"
	"path/filepath"
)

var fileIndex int

func writeTree(t *classification.Tree) {
	fileIndex += 1
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
		t := &classification.Tree{
			Hist: make([]int, p),
		}
		t.Hist[p/2] = p
		writeTree(t)
	}
	for p := 1; p < 600; p *= 2 {
		t := &classification.Tree{
			Hist: make([]int, p),
			LeftChild: &classification.Tree{
				Hist: make([]int, p),
			},
			RightChild: &classification.Tree{
				Hist: make([]int, p),
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
