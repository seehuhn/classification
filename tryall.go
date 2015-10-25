// +build ignore

package main

import (
	"fmt"
	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/stop"
	"github.com/seehuhn/classification/tree"
	"math"
	"strings"
)

type TreeFactory struct {
	name    string
	builder *tree.Factory
}

func (tf *TreeFactory) Name() string {
	return tf.name
}

func (tf *TreeFactory) FromTrainingData(data *classification.TrainingData) classification.Classifier {
	tree, _ := tf.builder.FromTrainingData(
		data.NumClasses, data.X, data.Y, data.Weight)
	return tree
}

func formatVal(x, se float64, width, maxPrec int) string {
	prec := int(-math.Log10(1.96 * se))
	if prec > maxPrec {
		prec = maxPrec
	}
	if prec < 0 {
		prec = 0
	}
	t := "%*.*f"
	t += strings.Repeat(" ", maxPrec-prec)
	return fmt.Sprintf(t, width+prec-maxPrec, prec, x)
}

type row struct {
	name   string
	values []<-chan *classification.Result
}

func main() {
	tree1 := &TreeFactory{"CART", tree.CART}
	tree2builder := &tree.Factory{
		StopGrowth: stop.IfPure,
		SplitScore: impurity.Gini,
		PruneScore: impurity.Gini,
		XValLoss:   loss.Deviance,
		K:          10,
	}
	tree2 := &TreeFactory{"other CART", tree2builder}
	methods := []classification.Factory{
		tree1,
		tree2,
	}

	testCases := []data.Set{
		data.Digits,
		data.NewNormals(1.0, 10000, 10000),
		data.NewNormals(2.0, 10000, 10000),
		data.NewNormals(3.0, 10000, 10000),
	}

	rows := make(chan *row, 8)
	go func() {
		for _, sample := range testCases {
			r := row{
				name:   sample.Name(),
				values: make([]<-chan *classification.Result, len(methods)),
			}
			for i, method := range methods {
				r.values[i] = classification.Assess(method, sample, loss.ZeroOne)
			}
			rows <- &r
		}
		close(rows)
	}()

	sampleNameLength := 25
	colWidth := 10
	maxPrec := 5

	fmt.Println()
	for i, method := range methods {
		fmt.Println(string([]byte{'A' + byte(i)}), "=", method.Name())
	}
	fmt.Println()
	fmt.Print(strings.Repeat(" ", sampleNameLength))
	for i := range methods {
		k := colWidth - maxPrec - 3
		fmt.Print("| " +
			strings.Repeat(" ", k) +
			string([]byte{'A' + byte(i)}) +
			strings.Repeat(" ", colWidth-k-3))
	}
	fmt.Println()
	fmt.Print(strings.Repeat("-", sampleNameLength))
	for range methods {
		fmt.Print("+" + strings.Repeat("-", colWidth-1))
	}
	fmt.Println()

	for row := range rows {
		fmt.Printf("%-*s", sampleNameLength, row.name)
		errors := []string{}
		for _, c := range row.values {
			value := <-c
			if value.Err != nil {
				fmt.Print("| ERROR" + strings.Repeat(" ", colWidth-7))
				errors = append(errors, value.Err.Error())
			} else {
				fmt.Print("| " + formatVal(value.MeanLoss, value.StdErr,
					colWidth-2, maxPrec))
			}
		}
		fmt.Println()
		for _, msg := range errors {
			fmt.Println("  " + msg)
		}
	}
	fmt.Println()
}
