package classification

import (
	"github.com/seehuhn/classification/util"
)

// A StopFunction is the type of function used to decide when to stop
// splitting branches of a classification tree, when the tree is
// originally constructed (i.e. before pruning).  The default stop
// function keeps adding branches until only one node is left.
type StopFunction func(util.Histogram) bool

func StopIfAtMost(n int) StopFunction {
	return func(freq util.Histogram) bool {
		return freq.Sum() <= n
	}
}
