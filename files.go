package classification

import (
	"fmt"
	"os"
)

func WriteMatrix(fname string, x *Matrix) {
	fd, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for i := 0; i < x.n; i++ {
		for j := 0; j < x.p; j++ {
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
