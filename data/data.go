package data

import (
	"github.com/seehuhn/classification/matrix"
)

// Data represents a collection of samples with associated classes,
// used either as trainings data or as test data in a classification
// problem.
type Data struct {
	// NumClasses gives the number of classes used in the data set.
	NumClasses int

	// X stores the predictor variables.  Each column of the matrix
	// corresponds to a variable, rows represent samples.
	X *matrix.Float64

	// Y stores the response variables.  The length of Y must equal
	// the number of rows of X.
	Y []int

	// Weights is either nil or a vector of the same length as Y.  If
	// the vector is present, it must the the same length as Y and
	// stores weights for each sample.  If `Weight` is nil, each
	// sample has weight 1.x
	Weights []float64

	// Rows can be either nil or a vector or row numbers.  If `Rows`
	// is non-nil, the entries must be valid row numbers for `X` and
	// the data set consists of the listed rows.  If `Rows` is nil,
	// the data set consists of all rows of `X`.
	Rows []int
}

// Get the vector of all rows used in the data set.  The returned
// slice is owned by the Data object and must not be changed by the
// caller.
func (data *Data) GetRows() []int {
	if data.Rows != nil {
		return data.Rows
	}
	res := make([]int, len(data.Y))
	for i := range res {
		res[i] = i
	}
	return res
}
