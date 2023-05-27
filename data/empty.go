package data

import (
	"seehuhn.de/go/classification/matrix"
)

// NewEmpty returns a newly allocated data set of the given size,
// where all variable values are zero.  The caller can change the
// entries of .X and .Y in the returned data set to fill in the data.
func NewEmpty(numClasses, n, p int) *Data {
	return &Data{
		NumClasses: numClasses,
		X:          matrix.NewFloat64(n, p, 0, nil),
		Y:          make([]int, n),
	}
}
