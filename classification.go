package classification

import (
	"math"
	"time"

	"seehuhn.de/go/classification/data"
	"seehuhn.de/go/classification/loss"
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

	TrainingTime time.Duration
	TestTime     time.Duration
}

// Assess estimates the quality of a classifier by computing the
// average loss, using Monte Carlo integration.  `samples` is used to
// construct training and test data, `cf` is used to construct the
// classifier from training data, and `L` specifies the loss function
// to assess the cost of wrong classifications.
func Assess(cf Factory, samples data.Set, L loss.Function) *Result {
	res := &Result{}

	trainingData, err := samples.TrainingData()
	if err != nil {
		res.Err = err
		return res
	}
	start := time.Now()
	c := cf.FromData(trainingData)
	res.TrainingTime = time.Since(start)

	testData, err := samples.TestData()
	if err != nil {
		res.Err = err
		return res
	}
	cumLoss := 0.0
	cumLoss2 := 0.0
	rows := testData.GetRows()
	start = time.Now()
	for _, i := range rows {
		sample := testData.X.Row(i)
		prob := c.EstimateClassProbabilities(sample)
		l := L(testData.Y[i], prob)
		cumLoss += l
		cumLoss2 += l * l
	}
	res.TestTime = time.Since(start)
	nn := float64(len(rows))
	cumLoss /= nn
	cumLoss2 /= nn
	stdErr := math.Sqrt((cumLoss2 - cumLoss*cumLoss) / (nn - 1))

	res.MeanLoss = cumLoss
	res.StdErr = stdErr
	return res
}
