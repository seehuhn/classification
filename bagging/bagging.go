package bagging

import (
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/data"
	"math/rand"
	"runtime"
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
	voterSize := f.VoterSize
	if voterSize == 0 {
		voterSize = data.NRow()
	}

	numWorkers := runtime.NumCPU()
	jobs := make(chan int64, f.NumVoters)
	for i := 0; i < f.NumVoters; i++ {
		jobs <- baggingSeed + int64(i)
	}
	close(jobs)

	results := make(chan classification.Classifier)
	for j := 0; j < numWorkers; j++ {
		go func() {
			for seed := range jobs {
				rng := rand.New(rand.NewSource(seed))
				sample := data.SampleWithReplacement(rng, voterSize)
				results <- f.Base.FromData(sample)
			}
		}()
	}

	res := make(baggingClassifier, f.NumVoters)
	for i := range res {
		res[i] = <-results
	}
	return res
}

type baggingClassifier []classification.Classifier

func (bag baggingClassifier) EstimateClassProbabilities(x []float64) data.Histogram {
	var res data.Histogram
	for _, cfr := range bag {
		p := cfr.EstimateClassProbabilities(x)
		if res == nil {
			res = make(data.Histogram, len(p))
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
