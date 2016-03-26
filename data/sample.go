package data

import (
	"math/rand"
)

// SampleWithoutReplacement returns a random subset of size `m` of the
// data.  The subset returned is chosen uniformly amongst all subsets
// of size m, the rows are returned in random order.  The method
// panics if the data set has fewer than m rows.
func (data *Data) SampleWithoutReplacement(rng *rand.Rand, m int) *Data {
	rows := data.GetRows()
	n := len(rows)
	if m > n {
		panic("requested sample size too large")
	}

	// use reservoir sampling:
	newRows := make([]int, m)
	copy(newRows, rows[:m])
	for i := m; i < n; i++ {
		j := rng.Intn(i + 1)
		if j < m {
			newRows[j] = rows[i]
		}
	}

	res := *data // make a shallow copy
	res.Rows = newRows
	return &res
}

// SampleWithReplacement returns a random subset of size `m` of the
// data.  Each element of the result is chosen uniformly amongst the
// elements of the data set, independently.
func (data *Data) SampleWithReplacement(rng *rand.Rand, m int) *Data {
	rows := data.GetRows()
	n := len(rows)

	newRows := make([]int, m)
	for i := range newRows {
		newRows[i] = rows[rng.Intn(n)]
	}

	res := *data // make a shallow copy
	res.Rows = newRows
	return &res
}
