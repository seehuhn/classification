package forest

import (
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/bagging"
)

type RandomForestFactory struct {
	RandomTree
	NumTrees int
}

func (f *RandomForestFactory) New() classification.Factory {
	return bagging.NewFromRandom(&f.RandomTree, f.NumTrees)
}
