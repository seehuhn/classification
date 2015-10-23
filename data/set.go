package data

import (
	"github.com/seehuhn/classification/matrix"
)

type Set interface {
	Name() string
	NumClasses() int
	TrainingSet() (X *matrix.Float64, Y []int, err error)
	TestSet() (X *matrix.Float64, Y []int, err error)
}
