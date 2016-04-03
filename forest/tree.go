package forest

import (
	"fmt"
	"math/rand"

	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/tree"
)

type RandomTree struct {
	NumSamples float64
	NumLeaves  int
	NumColumns int // number of columns to use for each split
	SplitScore impurity.Function
}

type leaf struct {
	node *tree.Tree
	d    *data.Data
}

func (f *RandomTree) GetName() string {
	return fmt.Sprintf("random tree %g/%d/%d",
		f.NumSamples, f.NumLeaves, f.NumColumns)
}

func (f *RandomTree) FromDataRandom(d *data.Data, rng *rand.Rand) classification.Classifier {
	numSamples := int(float64(d.NRow()) * f.NumSamples)
	sample := d.SampleWithoutReplacement(numSamples, rng)
	root := &tree.Tree{
		Hist: sample.GetHist(),
	}
	todo := f.NumLeaves - 1

	current := []leaf{
		{root, sample},
	}
	var next []leaf
	for todo > 0 && len(current) > 0 {
		i := rng.Intn(len(current))
		this := current[i]
		current = append(current[:i], current[i+1:]...)

		best := f.findBestSplit(rng, this.d, this.node.Hist)

		leftChild := &tree.Tree{
			Hist: best.LeftHist,
		}
		rightChild := &tree.Tree{
			Hist: best.RightHist,
		}
		this.node.LeftChild = leftChild
		this.node.RightChild = rightChild
		this.node.Column = best.Col
		this.node.Limit = best.Limit
		todo--

		if best.Left.NRow() > 1 {
			next = append(next, leaf{leftChild, best.Left})
		}
		if best.Right.NRow() > 1 {
			next = append(next, leaf{rightChild, best.Right})
		}

		if len(current) == 0 {
			current = next
			next = nil
		}
	}

	return root
}
