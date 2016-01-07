package bagging

import (
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/util"
	"math/rand"
)

const baggingSeed = 1070630982

type baggingFactory struct {
	Base      classification.Factory
	NumVoters int
	VoterSize int
}

func (f *baggingFactory) Name() string {
	if f.VoterSize == 0 {
		return fmt.Sprintf("%s, %d-bagged", f.Base.Name(), f.NumVoters)
	}
	return fmt.Sprintf("%s, %dx%d-bagged", f.Base.Name(), f.NumVoters, f.VoterSize)
}

// New constructs a new `classification.Factory`, constructed from the
// `base` classifier using bagging.  The resulting classifier
// aggregates the output of `numVoters` individual classifiers, each
// of which is trained using `voterSize` training samples.
func New(base classification.Factory, numVoters, voterSize int) classification.Factory {
	return &baggingFactory{
		Base:      base,
		NumVoters: numVoters,
		VoterSize: voterSize,
	}
}

func (f *baggingFactory) FromData(data *data.Data) classification.Classifier {
	rng := rand.New(rand.NewSource(baggingSeed))
	n := len(data.Y)
	voterSize := f.VoterSize
	if voterSize == 0 {
		voterSize = n
	}

	res := make(baggingClassifier, f.NumVoters)
	newData := *data // make a shallow copy
	newData.Rows = make([]int, voterSize)
	for j := range res {
		for i := range newData.Rows {
			newData.Rows[i] = rng.Intn(n)
		}
		res[j] = f.Base.FromData(&newData)
	}

	return res
}

type baggingClassifier []classification.Classifier

func (bag baggingClassifier) EstimateClassProbabilities(x []float64) util.Histogram {
	var res util.Histogram
	for _, cfr := range bag {
		p := cfr.EstimateClassProbabilities(x)
		if res == nil {
			res = make(util.Histogram, len(p))
		}
		for i, pi := range p {
			res[i] += pi
		}
	}
	for i := range res {
		res[i] /= float64(len(bag))
	}
	return res
}
