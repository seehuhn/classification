package classification

func probabilities(freq []int) []float64 {
	prob := make([]float64, len(freq))
	n := sum(freq)
	for i, ni := range freq {
		prob[i] = float64(ni) / float64(n)
	}
	return prob
}

func sum(freq []int) int {
	res := 0
	for _, ni := range freq {
		res += ni
	}
	return res
}
