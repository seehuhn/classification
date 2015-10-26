package data

import (
	"github.com/seehuhn/classification/matrix"
)

type Data struct {
	NumClasses int
	X          *matrix.Float64
	Y          []int
	Weights    []float64
	Rows       []int
}

type Set interface {
	Name() string
	TrainingData() (data *Data, err error)
	TestData() (data *Data, err error)
}
