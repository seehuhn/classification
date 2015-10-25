package classification

import (
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/util"
	"math"
)

type Classifier interface {
	EstimateClassProbabilities(x []float64) util.Histogram
}

type Factory interface {
	Name() string
	FromData(*data.Data) Classifier
}

type Result struct {
	MeanLoss float64
	StdErr   float64
	Err      error
}

func Assess(cf Factory, samples data.Set, L loss.Function) *Result {
	trainingData, err := samples.TrainingData()
	if err != nil {
		return &Result{0, 0, err}
	}
	c := cf.FromData(trainingData)

	testData, err := samples.TestData()
	if err != nil {
		return &Result{0, 0, err}
	}

	cumLoss := 0.0
	cumLoss2 := 0.0
	nTest := len(testData.Y)
	for i := 0; i < nTest; i++ {
		row := testData.X.Row(i)
		prob := c.EstimateClassProbabilities(row)
		l := L(testData.Y[i], prob)
		cumLoss += l
		cumLoss2 += l * l
	}
	nn := float64(nTest)
	cumLoss /= nn
	cumLoss2 /= nn
	stdErr := math.Sqrt((cumLoss2 - cumLoss*cumLoss) / nn)

	return &Result{cumLoss, stdErr, nil}
}
