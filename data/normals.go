package data

import (
	"fmt"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/mt19937"
	"math/rand"
)

type twoNormals struct {
	name      string
	X         *matrix.Float64
	Y         []int
	trainRows []int
	testRows  []int
}

// NewNormals is a factory function which can generate data sets.  The
// resulting data sets consist of mixtures of two one-dimensional
// normal distributions, with variance 1, where the means are
// separated by `delta`.
func NewNormals(delta float64, nTrain, nTest int) Set {
	nTotal := nTrain + nTest

	X := make([]float64, nTotal)
	Y := make([]int, nTotal)
	rows := make([]int, nTotal)

	rng := rand.New(mt19937.New())
	means := []float64{-delta / 2, +delta / 2}
	for i := range rows {
		class := rng.Intn(2)
		Y[i] = class
		X[i] = rng.NormFloat64() + means[class]
		rows[i] = i
	}

	return &twoNormals{
		name:      fmt.Sprintf("normals/%g/%d/%d", delta, nTrain, nTest),
		X:         matrix.NewFloat64(nTotal, 1, 0, X),
		Y:         Y,
		trainRows: rows[:nTrain],
		testRows:  rows[nTrain:],
	}
}

func (ss *twoNormals) Name() string {
	return ss.name
}

func (ss *twoNormals) TrainingData() (data *Data, err error) {
	return &Data{
		NumClasses: 2,
		X:          ss.X,
		Y:          ss.Y,
		Rows:       ss.trainRows,
	}, nil
}

func (ss *twoNormals) TestData() (data *Data, err error) {
	return &Data{
		NumClasses: 2,
		X:          ss.X,
		Y:          ss.Y,
		Rows:       ss.testRows,
	}, nil
}
