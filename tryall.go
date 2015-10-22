// +build ignore

package main

import (
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/matrix"
	"github.com/seehuhn/classification/tree"
)

type TreeFactory struct {
	name    string
	builder *tree.Builder
}

func (tf *TreeFactory) Name() string {
	return tf.name
}

func (tf *TreeFactory) FromTrainingData(numClasses int, X *matrix.Float64,
	Y []int, weight []float64) classification.Classifier {
	tree, _ := tf.builder.NewFromTrainingData(numClasses, X, Y, weight)
	return tree
}

func main() {
	fmt.Println("hello")

	tree := &TreeFactory{"CART", tree.CART}
	samples := classification.NewTwoNormals(1.0)
	classification.Assess(tree, samples, 1000, 10000)
	samples = classification.NewTwoNormals(2.0)
	classification.Assess(tree, samples, 1000, 10000)
}
