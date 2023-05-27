package stop

import (
	"seehuhn.de/go/classification/data"
)

// A Function is used to decide when to stop splitting branches of a
// classification tree, when the tree is originally constructed
// (i.e. before pruning).  The default stop function keeps adding
// branches until only one node is left.
type Function func(data.Histogram) bool

// TODO(voss): use different naming conventions for stop functions and
// stop function factories?

// IfAtMost returns a stop function which stops splitting nodes once
// the current node has n or fewer samples associated to it.
func IfAtMost(n float64) Function {
	return func(hist data.Histogram) bool {
		return hist.Sum() <= n
	}
}

// IfPure is a [Function] which stops splitting nodes once
// all samples in the current node have the same class.
func IfPure(hist data.Histogram) bool {
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

// IfPureOrAtMost returns a stop function which stops splitting nodes
// when either all samples in the current node have the same class, or
// the node has `n` or fewer samples associated to it.
func IfPureOrAtMost(n float64) Function {
	return func(hist data.Histogram) bool {
		total := 0.0
		k := 0
		for _, ni := range hist {
			total += ni
			if ni > 0 {
				k++
			}
			if k > 1 && total > n {
				return false
			}
		}
		return true
	}
}
