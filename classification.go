package classification

import (
	"fmt"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"time"
)

type Classifier interface {
	EstimateClassProbabilities(x []float64) []float64
	GuessClass(x []float64) int
}

type Factory interface {
	Name() string
	FromTrainingData(
		numClasses int,
		X *matrix.Float64,
		Y []int,
		weight []float64) Classifier
}

type SimulatedSamples interface {
	Name() string
	NumClasses() int
	NewSamples(n int) (X *matrix.Float64, Y []int)
}

type twoNormals struct {
	means []float64
	rng   *rand.Rand
}

func NewTwoNormals(delta float64) SimulatedSamples {
	res := &twoNormals{
		means: []float64{0.0, delta},
		rng:   rand.New(mt19937.New()),
	}
	res.rng.Seed(time.Now().UnixNano())
	return res
}

func (ss *twoNormals) NumClasses() int {
	return 2
}

func (ss *twoNormals) NewSamples(n int) (X *matrix.Float64, Y []int) {
	X = matrix.NewFloat64(n, 1, 0, nil)
	Y = make([]int, n)
	for i := 0; i < n; i++ {
		Yi := ss.rng.Intn(2)
		Y[i] = Yi
		X.Set(i, 0, ss.rng.NormFloat64()+ss.means[Yi])
	}
	return
}

func (ss *twoNormals) Name() string {
	return fmt.Sprintf("two normals (delta=%g)", ss.means[1])
}

func Assess(classifier Factory, samples SimulatedSamples, nTrain, nTest int) {
	numClasses := samples.NumClasses()
	XTrain, YTrain := samples.NewSamples(nTrain)
	c := classifier.FromTrainingData(numClasses, XTrain, YTrain, nil)

	res := matrix.NewInt(numClasses, numClasses, 0, nil)
	XTest, YTest := samples.NewSamples(nTest)
	for i := 0; i < nTest; i++ {
		row := XTest.Row(i)
		guessed := c.GuessClass(row)
		correct := YTest[i]
		res.Inc(correct, guessed)
	}
	fmt.Println(res)
}
