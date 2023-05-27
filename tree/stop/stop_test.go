package stop

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type Tests struct{}

var _ = Suite(&Tests{})

func (*Tests) TestIfAtMost(c *C) {
	atMost7 := IfAtMost(7)
	c.Assert(atMost7([]float64{1, 1, 1, 1, 1, 1, 1}), Equals, true)
	c.Assert(atMost7([]float64{3, 0, 4}), Equals, true)
	c.Assert(atMost7([]float64{7, 0}), Equals, true)
	c.Assert(atMost7([]float64{7}), Equals, true)
	c.Assert(atMost7([]float64{1}), Equals, true)

	c.Assert(atMost7([]float64{1, 1, 1, 1, 1, 1, 1, 1}), Equals, false)
	c.Assert(atMost7([]float64{3, 0, 5}), Equals, false)
	c.Assert(atMost7([]float64{8, 0}), Equals, false)
	c.Assert(atMost7([]float64{8}), Equals, false)
	c.Assert(atMost7([]float64{100}), Equals, false)
}

func (*Tests) TestPure(c *C) {
	c.Assert(IfPure([]float64{1, 0, 0}), Equals, true)
	c.Assert(IfPure([]float64{0, 1, 0}), Equals, true)
	c.Assert(IfPure([]float64{0, 0, 1}), Equals, true)
	c.Assert(IfPure([]float64{0, 100, 0}), Equals, true)
	c.Assert(IfPure([]float64{2}), Equals, true)

	c.Assert(IfPure([]float64{100, 0, 0, 1, 0}), Equals, false)
	c.Assert(IfPure([]float64{1, 1}), Equals, false)
	c.Assert(IfPure([]float64{1, 0, 1}), Equals, false)
	c.Assert(IfPure([]float64{0, 0, 1, 0, 1}), Equals, false)
	c.Assert(IfPure([]float64{2, 2}), Equals, false)
}

func (*Tests) TestIfPureOrAtMost(c *C) {
	stop := IfPureOrAtMost(7)
	c.Assert(stop([]float64{0, 0, 1}), Equals, true)
	c.Assert(stop([]float64{0, 1, 0}), Equals, true)
	c.Assert(stop([]float64{0, 100, 0}), Equals, true)
	c.Assert(stop([]float64{1, 0, 0}), Equals, true)
	c.Assert(stop([]float64{1, 1, 1, 1, 1, 1, 1}), Equals, true)
	c.Assert(stop([]float64{2}), Equals, true)
	c.Assert(stop([]float64{3, 0, 4}), Equals, true)
	c.Assert(stop([]float64{7, 0}), Equals, true)
	c.Assert(stop([]float64{7}), Equals, true)
	c.Assert(stop([]float64{8, 0}), Equals, true)

	c.Assert(stop([]float64{1, 1, 1, 1, 0, 1, 1, 1, 1}), Equals, false)
	c.Assert(stop([]float64{1, 1, 1, 1, 1, 1, 1, 1}), Equals, false)
	c.Assert(stop([]float64{3, 0, 5}), Equals, false)
	c.Assert(stop([]float64{100, 100}), Equals, false)
}
