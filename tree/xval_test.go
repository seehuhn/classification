package tree

import (
	. "gopkg.in/check.v1"
	"sort"
)

func (*Tests) TestXValClasses(c *C) {
	for _, n := range []int{5, 49, 50, 51} {
		for _, K := range []int{2, 3, 4, 5} {
			allTests := []int{}
			for k := 0; k < K; k++ {
				trainingSet, testSet := getXValSets(k, K, n)
				c.Check(len(trainingSet) >= len(testSet)-1, Equals, true)
				allTests = append(allTests, testSet...)

				together := append(trainingSet, testSet...)
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
