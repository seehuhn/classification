package data

import (
	"fmt"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"sync"
)

type twoNormals struct {
	sync.Mutex
	means  []float64
	rng    *rand.Rand
	nTrain int
	nTest  int
}

func NewNormals(delta float64, nTrain, nTest int) Set {
	res := &twoNormals{
		means:  []float64{0.0, delta},
		rng:    rand.New(mt19937.New()),
		nTrain: nTrain,
		nTest:  nTest,
	}
	return res
}

func (ss *twoNormals) Name() string {
	return fmt.Sprintf("normals/%g/%d/%d", ss.means[1], ss.nTrain, ss.nTest)
}

func (ss *twoNormals) NumClasses() int {
	return 2
}

func (ss *twoNormals) newSamples(n int) (X *matrix.Float64, Y []int, err error) {
	ss.Lock()
	defer ss.Unlock()
	X = matrix.NewFloat64(n, 1, 0, nil)
	Y = make([]int, n)
	for i := 0; i < n; i++ {
		Yi := ss.rng.Intn(2)
		Y[i] = Yi
		X.Set(i, 0, ss.rng.NormFloat64()+ss.means[Yi])
	}
	return
}

func (ss *twoNormals) TrainingSet() (X *matrix.Float64, Y []int, err error) {
	return ss.newSamples(ss.nTrain)
}

func (ss *twoNormals) TestSet() (X *matrix.Float64, Y []int, err error) {
	return ss.newSamples(ss.nTest)
}
