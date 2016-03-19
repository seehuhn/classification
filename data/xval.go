package data

import (
	"fmt"
	"math/rand"
)

type xvalSet struct {
	k, K        int
	data        *Data
	trainingSet []int
	testSet     []int
}

func (xv *xvalSet) Name() string {
	return fmt.Sprintf("XvalSet %d/%d", xv.k, xv.K)
}

func (xv *xvalSet) TrainingData() (data *Data, err error) {
	res := *xv.data
	res.Rows = xv.trainingSet
	return &res, nil
}

func (xv *xvalSet) TestData() (data *Data, err error) {
	res := *xv.data
	res.Rows = xv.testSet
	return &res, nil
}

func (data *Data) GetXValSet(seed int64, k, K int) Set {
	if K < 2 {
		panic("need at least K=2 groups for cross-validation")
	}

	rows := data.GetRows()
	n := len(rows)
	if n < K {
		panic("not enough samples for cross-validation")
	}

	rng := rand.New(rand.NewSource(seed))
	shuffled := make([]int, n)
	// A simplified version of
	// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_.22inside-out.22_algorithm
	for i, row := range rows {
		// Potential optimisations: for i = j, the assignment
		// shuffled[i] = shuffled[j] is unnecessary.  For i = 0 we
		// don't even need to all Intn().
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], row
	}

	a := k * n / K
	b := (k + 1) * n / K

	// This could be optimised if needed: for example, for k=0 and
	// k=K-1 we don't need to allocate a new slice.
	testSet := make([]int, b-a)
	copy(testSet, shuffled[a:b])
	trainingSet := append(shuffled[:a], shuffled[b:]...)

	res := &xvalSet{
		k:           k,
		K:           K,
		data:        data,
		trainingSet: trainingSet,
		testSet:     testSet,
	}
	return res
}
