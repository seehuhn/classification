package classification

import (
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/util"
	"math"
	"runtime"
)

type Classifier interface {
	EstimateClassProbabilities(x []float64) util.Histogram
}

type Factory interface {
	Name() string
	FromTrainingData(
		numClasses int,
		X *matrix.Float64,
		Y []int,
		weight []float64) Classifier
}

type Result struct {
	MeanLoss float64
	StdErr   float64
	Err      error
}

func doAssess(cf Factory, samples data.Set, L loss.Function) *Result {
	numClasses := samples.NumClasses()
	XTrain, YTrain, err := samples.TrainingSet()
	if err != nil {
		return &Result{0, 0, err}
	}
	c := cf.FromTrainingData(numClasses, XTrain, YTrain, nil)

	XTest, YTest, err := samples.TestSet()
	if err != nil {
		return &Result{0, 0, err}
	}

	cumLoss := 0.0
	cumLoss2 := 0.0
	nTest := len(YTest)
	for i := 0; i < nTest; i++ {
		row := XTest.Row(i)
		prob := c.EstimateClassProbabilities(row)
		l := L(YTest[i], prob)
		cumLoss += l
		cumLoss2 += l * l
	}
	nn := float64(nTest)
	cumLoss /= nn
	cumLoss2 /= nn
	stdErr := math.Sqrt((cumLoss2 - cumLoss*cumLoss) / nn)

	return &Result{cumLoss, stdErr, nil}
}

var queue chan int

func Assess(cf Factory, samples data.Set, L loss.Function) <-chan *Result {
	worker := <-queue
	resChan := make(chan *Result, 1)
	go func() {
		res := doAssess(cf, samples, L)
		resChan <- res
		queue <- worker
	}()
	return resChan
}

func init() {
	n := runtime.GOMAXPROCS(0)
	queue = make(chan int, n)
	for i := 0; i < n; i++ {
		queue <- i
	}
}
