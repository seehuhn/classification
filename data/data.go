package data

import (
	"github.com/seehuhn/classification/matrix"
)

// Data represents a collection of samples with associated classes,
// used either as trainings data or as test data.
type Data struct {
	NumClasses int
	X          *matrix.Float64
	Y          []int
	Weights    []float64
	Rows       []int
}

// Set implements an abstract interface to represent a test data set,
// consisting of trainings data for setting up the method and test
// data for assessment.
type Set interface {
	Name() string
	TrainingData() (data *Data, err error)
	TestData() (data *Data, err error)
}
