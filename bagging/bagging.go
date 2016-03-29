package bagging

import (
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/data"
	"math/rand"
	"runtime"
)

const baggingSeed = 1070630982

type RandomFactory interface {
	Name() string
	FromDataRandom(d *data.Data, rng *rand.Rand) classification.Classifier
}

type randomize struct {
	base      classification.Factory
	voterSize int
}

func (f randomize) Name() string {
	return "random " + f.base.Name()
}

func (f randomize) FromDataRandom(d *data.Data, rng *rand.Rand) classification.Classifier {
	voterSize := f.voterSize
	if voterSize == 0 {
		voterSize = d.NRow()
	}

	sample := d.SampleWithReplacement(voterSize, rng)
	return f.base.FromData(sample)
}

// New constructs a new `classification.Factory`, using the `base`
// classifier together with bagging.  The resulting classifier
// aggregates the output of `numVoters` individual classifiers, each
// of which is trained using `voterSize` training samples.
func New(base classification.Factory, numVoters, voterSize int) classification.Factory {
	return &baggingFactory{
		Base:      randomize{base, voterSize},
		NumVoters: numVoters,
	}
}

func NewFromRandom(base RandomFactory, numVoters int) classification.Factory {
	return &baggingFactory{
		Base:      base,
		NumVoters: numVoters,
	}
}

type baggingFactory struct {
	Base      RandomFactory
	NumVoters int
}

func (f *baggingFactory) Name() string {
	return fmt.Sprintf("%s, %d-bagged", f.Base.Name(), f.NumVoters)
}

func (f *baggingFactory) FromData(data *data.Data) classification.Classifier {
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
				results <- f.Base.FromDataRandom(data, rng)
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
	// TODO(voss): should averages take leaf size into account for trees?
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
