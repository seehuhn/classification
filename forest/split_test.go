package forest

import (
	. "gopkg.in/check.v1"
	"math/rand"
)

func (*Tests) TestSubset(c *C) {
	rngSource := rand.NewSource(0)
	rng := rand.New(rngSource)

	N := 1000
	m := 3
	n := 5
	count := make([]uint, n)
	for i := 0; i < N; i++ {
		subset := subset(rng, m, n)
		c.Assert(len(subset), Equals, m)
		c.Assert(subset[0], Not(Equals), subset[1])
		for j := 0; j < m; j++ {
			count[subset[j]]++
		}
	}

	// Use a chi-squared test.
	expected := float64(m) * float64(N) / float64(n)
	chiSquared := 0.0
	for i := 0; i < n; i++ {
		delta := float64(count[i]) - expected
		chiSquared += delta * delta / expected
	}
	// Four degress of freedom, alpha = 0.01 (i.e. only one
	// percent of all seeds should fail this test if the algorithm
	// is implemented correctly).
	limit := 13.28
	if chiSquared > limit {
		c.Errorf("chi-squared test failed, %f > %f", chiSquared, limit)
	}
}
