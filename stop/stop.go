package stop

import (
	"github.com/seehuhn/classification/util"
)

// A Function is used to decide when to stop splitting branches of a
// classification tree, when the tree is originally constructed
// (i.e. before pruning).  The default stop function keeps adding
// branches until only one node is left.
type Function func(util.Histogram) bool

// TODO(voss): use different naming conventions for stop functions and
// stop function factories.

// IfAtMost returns a stop.Function which stops splitting nodes once
// the current node has `n` or fewer samples associated to it.
func IfAtMost(n int) Function {
	return func(freq util.Histogram) bool {
		return freq.Sum() <= n
	}
}

// IfHomogeneous is a stop.Function which stops splitting nodes once
// all samples in the current node have the same class.
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
