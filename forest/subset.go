package forest

import (
	"math/rand"
)

// Subset returns a random subset of {0, 1, ..., n-1} with m elements.
func Subset(r *rand.Rand, m, n int) []int {
	if m > n {
		panic("invalid subset size")
	}

	// use reservoir sampling:
	res := make([]int, m)
	for i := 0; i < m; i++ {
		res[i] = i
	}
	for i := m; i < n; i++ {
		j := r.Intn(i + 1)
		if j < m {
			res[j] = i
		}
	}
	return res
}
