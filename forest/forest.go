package forest

import (
	"seehuhn.de/go/classification"
	"seehuhn.de/go/classification/bagging"
)

type RandomForestFactory struct {
	RandomTree
	NumTrees int
}

func (f *RandomForestFactory) New() classification.Factory {
	return bagging.NewFromRandom(&f.RandomTree, f.NumTrees)
}
