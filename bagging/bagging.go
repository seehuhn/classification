package bagging

import (
	"github.com/seehuhn/classification"
)

type baggingClassifier struct {
	voters []classification.Classifier
}

type baggingFactory struct {
	in classification.Factory
	n  int
}

func (f *baggingFactory) Name() string {
	return f.in.Name + " (bagged)"
}

func (f *baggingFactory) FromTrainingData(numClasses int, X *matrix.Float64,
	Y []int, weight []float64) classification.Classifier {
	res := &baggingClassifier{}
	res.voters = make([]classification.Classifier, f.n)
	for i := 0; i < f.n; i++ {

	}
	return res
}

func New(in classification.Factory, n int) classification.Factory {
	return &baggingFactory{
		in: in,
		n:  n,
	}
}
