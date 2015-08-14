package stop

import (
	"github.com/seehuhn/classification/util"
)

// A stop.Function is used to decide when to stop splitting branches
// of a classification tree, when the tree is originally constructed
// (i.e. before pruning).  The default stop function keeps adding
// branches until only one node is left.
type Function func(util.Histogram) bool

// TODO(voss): use different naming conventions for stop functions and
// factories.

func IfAtMost(n int) Function {
	return func(freq util.Histogram) bool {
		return freq.Sum() <= n
	}
}

func IfHomogeneous(hist util.Histogram) bool {
	k := 0
	for _, ni := range hist {
		if ni > 0 {
			k++
			if k > 1 {
				return false
			}
		}
	}
	return true
}
