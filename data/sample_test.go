package data

import (
	"math/rand"

	. "gopkg.in/check.v1"
)

func (*Tests) TestSubset(c *C) {
	rngSource := rand.NewSource(0)
	rng := rand.New(rngSource)

	N := 1000
	m := 3
	n := 5
	data := NewEmpty(2, n, 0)
	count := make([]uint, n)
	for i := 0; i < N; i++ {
		subset := data.SampleWithoutReplacement(m, rng)
		rows := subset.GetRows()
		c.Assert(len(rows), Equals, m)
		c.Assert(rows[0], Not(Equals), rows[1])
		for j := 0; j < m; j++ {
			count[rows[j]]++
		}
	}

	// Use a chi-squared test.
	expected := float64(m) * float64(N) / float64(n)
	chiSquared := 0.0
	for i := 0; i < n; i++ {
		delta := float64(count[i]) - expected
		chiSquared += delta * delta / expected
	}
	// We have n-1=4 degress of freedom and choose alpha = 0.01 (i.e.,
	// if the algorithm is implemented correctly, this test should
	// fail only for one percent of all seeds).
	limit := 13.27670
	if chiSquared > limit {
		c.Errorf("chi-squared test failed, %f > %f", chiSquared, limit)
	}
}
