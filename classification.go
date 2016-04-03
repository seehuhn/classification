package classification

import (
	"math"

	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/loss"
)

// Classifier represents an algorithm for classification, together
// with the information contributed by a training data set.
type Classifier interface {
	EstimateClassProbabilities(x []float64) data.Histogram
}

// Factory objects encapsulate the logic to construct classifiers from
// training data sets.
type Factory interface {
	GetName() string
	FromData(*data.Data) Classifier
}

// Result is returned by the `Assess` function to describe the quality
// of a classifier.
type Result struct {
	MeanLoss float64
	StdErr   float64
	Err      error
}

// Assess estimates the quality of a classifier by computing the
// average loss, using Monte Carlo integration.  `samples` is used to
// construct training and test data, `cf` is used to construct the
// classifier from training data, and `L` specifies the loss function
// to assess the cost of wrong classifications.
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
	rows := testData.GetRows()
	for _, i := range rows {
		sample := testData.X.Row(i)
		prob := c.EstimateClassProbabilities(sample)
		l := L(testData.Y[i], prob)
		cumLoss += l
		cumLoss2 += l * l
	}
	nn := float64(len(rows))
	cumLoss /= nn
	cumLoss2 /= nn
	stdErr := math.Sqrt((cumLoss2 - cumLoss*cumLoss) / (nn - 1))

	return &Result{cumLoss, stdErr, nil}
}
