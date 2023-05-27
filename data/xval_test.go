package data

import (
	"sort"

	. "gopkg.in/check.v1"
	"seehuhn.de/go/classification/matrix"
)

func dummyData(n int) *Data {
	return &Data{
		NumClasses: 2,
		X:          matrix.NewFloat64(n, 0, 0, nil),
		Y:          make([]int, n),
	}
}

func (*Tests) TestXValClasses(c *C) {
	seed := int64(1234567890)
	for _, n := range []int{5, 49, 50, 51} {
		data := dummyData(n)
		for _, K := range []int{2, 3, 4, 5} {
			var allTests []int
			for k := 0; k < K; k++ {
				set := data.GetXValSet(seed, K, k)
				td, err := set.TrainingData()
				c.Assert(err, IsNil)
				trainingRows := td.GetRows()
				td, err = set.TestData()
				c.Assert(err, IsNil)
				testRows := td.GetRows()
				c.Check(len(trainingRows) >= len(testRows)-1, Equals, true)
				c.Check(len(testRows) >= n/K, Equals, true)
				c.Check(len(testRows) <= n/K+1, Equals, true)
				allTests = append(allTests, testRows...)

				together := append(trainingRows, testRows...)
				c.Check(len(together), Equals, n)
				sort.IntSlice(together).Sort()
				for i := 0; i < n; i++ {
					c.Check(together[i], Equals, i)
				}
			}
			c.Check(allTests, HasLen, n)
			sort.IntSlice(allTests).Sort()
			for i := 0; i < n; i++ {
				c.Check(allTests[i], Equals, i)
			}
		}
	}
}
