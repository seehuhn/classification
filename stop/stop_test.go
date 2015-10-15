package stop

import (
	"github.com/seehuhn/classification/util"
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type Tests struct{}

var _ = Suite(&Tests{})

func (*Tests) TestIfAtMost(c *C) {
	atMost7 := IfAtMost(7)
	c.Assert(atMost7(util.Histogram{1, 1, 1, 1, 1, 1, 1}), Equals, true)
	c.Assert(atMost7(util.Histogram{3, 0, 4}), Equals, true)
	c.Assert(atMost7(util.Histogram{7, 0}), Equals, true)
	c.Assert(atMost7(util.Histogram{7}), Equals, true)
	c.Assert(atMost7(util.Histogram{1}), Equals, true)

	c.Assert(atMost7(util.Histogram{1, 1, 1, 1, 1, 1, 1, 1}), Equals, false)
	c.Assert(atMost7(util.Histogram{3, 0, 5}), Equals, false)
	c.Assert(atMost7(util.Histogram{8, 0}), Equals, false)
	c.Assert(atMost7(util.Histogram{8}), Equals, false)
	c.Assert(atMost7(util.Histogram{100}), Equals, false)
}

func (*Tests) TestPure(c *C) {
	c.Assert(IfPure(util.Histogram{1, 0, 0}), Equals, true)
	c.Assert(IfPure(util.Histogram{0, 1, 0}), Equals, true)
	c.Assert(IfPure(util.Histogram{0, 0, 1}), Equals, true)
	c.Assert(IfPure(util.Histogram{0, 100, 0}), Equals, true)
	c.Assert(IfPure(util.Histogram{2}), Equals, true)

	c.Assert(IfPure(util.Histogram{100, 0, 0, 1, 0}), Equals, false)
	c.Assert(IfPure(util.Histogram{1, 1}), Equals, false)
	c.Assert(IfPure(util.Histogram{1, 0, 1}), Equals, false)
	c.Assert(IfPure(util.Histogram{0, 0, 1, 0, 1}), Equals, false)
	c.Assert(IfPure(util.Histogram{2, 2}), Equals, false)
}

func (*Tests) TestIfPureOrAtMost(c *C) {
	stop := IfPureOrAtMost(7)
	c.Assert(stop(util.Histogram{0, 0, 1}), Equals, true)
	c.Assert(stop(util.Histogram{0, 1, 0}), Equals, true)
	c.Assert(stop(util.Histogram{0, 100, 0}), Equals, true)
	c.Assert(stop(util.Histogram{1, 0, 0}), Equals, true)
	c.Assert(stop(util.Histogram{1, 1, 1, 1, 1, 1, 1}), Equals, true)
	c.Assert(stop(util.Histogram{2}), Equals, true)
	c.Assert(stop(util.Histogram{3, 0, 4}), Equals, true)
	c.Assert(stop(util.Histogram{7, 0}), Equals, true)
	c.Assert(stop(util.Histogram{7}), Equals, true)
	c.Assert(stop(util.Histogram{8, 0}), Equals, true)

	c.Assert(stop(util.Histogram{1, 1, 1, 1, 0, 1, 1, 1, 1}), Equals, false)
	c.Assert(stop(util.Histogram{1, 1, 1, 1, 1, 1, 1, 1}), Equals, false)
	c.Assert(stop(util.Histogram{3, 0, 5}), Equals, false)
	c.Assert(stop(util.Histogram{100, 100}), Equals, false)
}
