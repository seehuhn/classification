package bagging

import (
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/util"
)

type baggingClassifier struct {
	voters []classification.Classifier
}

func (bag *baggingClassifier) EstimateClassProbabilities(x []float64) util.Histogram {
	// TODO(voss): implement this
	panic("not implemented")
	return nil
}

type baggingFactory struct {
	in classification.Factory
	n  int
}

func (f *baggingFactory) Name() string {
	return f.in.Name() + " (bagged)"
}

func (f *baggingFactory) FromTrainingData(data *classification.TrainingData) classification.Classifier {
	res := &baggingClassifier{}
	res.voters = make([]classification.Classifier, f.n)
	for i := 0; i < f.n; i++ {
		// TODO(voss): implement this
		panic("not implemented")
	}
	return res
}

func New(in classification.Factory, n int) classification.Factory {
	return &baggingFactory{
		in: in,
		n:  n,
	}
}
