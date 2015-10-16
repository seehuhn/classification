package tree

import (
	"fmt"
	. "gopkg.in/check.v1"
	"math"
)

type approxChecker struct {
	*CheckerInfo
}

var ApprixmatelyEquals Checker = &approxChecker{
	&CheckerInfo{
		Name:   "ApprixmatelyEquals",
		Params: []string{"obtained", "expected"},
	},
}

func (checker *approxChecker) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()
	a := math.Abs(params[0].(float64) - params[1].(float64))
	b := math.Abs(params[1].(float64))
	if b < 1e-6 {
		return a < 1e-6, ""
	}
	return a/b < 1e-6, ""
}

func (*Tests) TestTuneProfileGet(c *C) {
	p := &tuneProfile{
		breaks: []float64{0.0, 1.0, 2.0},
		values: []float64{5.0, 6.0, 7.0, 8.0},
	}
	c.Check(p.Get(-0.5), ApprixmatelyEquals, 5.0)
	c.Check(p.Get(0.5), ApprixmatelyEquals, 6.0)
	c.Check(p.Get(1.5), ApprixmatelyEquals, 7.0)
	c.Check(p.Get(2.5), ApprixmatelyEquals, 8.0)
}

func (*Tests) TestTuneProfileAdd(c *C) {
	p := &tuneProfile{}

	p.Add([]float64{}, []float64{1.0})
	c.Check(p.Get(0.0), ApprixmatelyEquals, 1.0)

	p.Add([]float64{0.0}, []float64{-1.0, 1.0})
	c.Check(p.Get(-1), ApprixmatelyEquals, 0.0)
	c.Check(p.Get(+1), ApprixmatelyEquals, 2.0)

	p.Add([]float64{-2.0, +2.0}, []float64{0.0, 10.0, 0.0})
	c.Check(p.Get(-3), ApprixmatelyEquals, 0.0)
	c.Check(p.Get(-1), ApprixmatelyEquals, 10.0)
	c.Check(p.Get(+1), ApprixmatelyEquals, 12.0)
	c.Check(p.Get(+3), ApprixmatelyEquals, 2.0)
}

func (*Tests) TestTuneProfileMinimum(c *C) {
	// left
	p := &tuneProfile{
		breaks: []float64{1.0, 2.0},
		values: []float64{5.0, 6.0, 7.0},
	}
	pos, val := p.Minimum()
	c.Check(pos, Equals, 0.0)
	c.Check(val, ApprixmatelyEquals, 5.0)

	// inside
	p = &tuneProfile{
		breaks: []float64{1.0, 2.0},
		values: []float64{5.0, 4.0, 7.0},
	}
	pos, val = p.Minimum()
	c.Check(pos, ApprixmatelyEquals, math.Sqrt(1.0*2.0))
	c.Check(val, ApprixmatelyEquals, 4.0)

	// right
	p = &tuneProfile{
		breaks: []float64{1.0, 2.0},
		values: []float64{5.0, 4.0, 3.0},
	}
	pos, val = p.Minimum()
	c.Check(pos, Equals, math.Inf(+1))
	c.Check(val, ApprixmatelyEquals, 3.0)

	// degenerate case
	p = &tuneProfile{
		breaks: []float64{},
		values: []float64{2.0},
	}
	_, val = p.Minimum()
	c.Check(val, ApprixmatelyEquals, 2.0)
}
