package classification

import (
	"fmt"
	"os"
)

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
