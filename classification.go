package classification

type Classifier interface {
	EstimateClassProbabilities(x []float64) []float64
	GuessClass(x []float64) int
}
