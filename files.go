package classification

import (
	"fmt"
	"github.com/gonum/matrix/mat64"
	"os"
)

func WriteMatrix(fname string, x mat64.Matrix) {
	fd, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	r, c := x.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if j > 0 {
				fd.WriteString(",")
			}
			fmt.Fprintf(fd, "%g", x.At(i, j))
		}
		fd.WriteString("\n")
	}
}

func WriteVector(fname string, x []int) {
	fd, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for _, xi := range x {
		fmt.Fprintf(fd, "%d\n", xi)
	}
}
