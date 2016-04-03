// +build ignore

package main

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/seehuhn/classification"
	"github.com/seehuhn/classification/bagging"
	"github.com/seehuhn/classification/data"
	"github.com/seehuhn/classification/forest"
	"github.com/seehuhn/classification/impurity"
	"github.com/seehuhn/classification/loss"
	"github.com/seehuhn/classification/stop"
	"github.com/seehuhn/classification/tree"
)

func formatVal(x, se float64, width, maxPrec int) string {
	prec := int(math.Ceil(-math.Log10(1.96 * se)))
	if prec > maxPrec {
		prec = maxPrec
	}
	if prec < 0 {
		prec = 0
	}
	t := "%*.*f" + strings.Repeat(" ", maxPrec-prec)
	return fmt.Sprintf(t, width+prec-maxPrec, prec, x)
}

type row struct {
	name   string
	values []<-chan *classification.Result
}

func main() {
	tree1 := tree.CART
	tree2 := &tree.Factory{
		Name:       "other CART",
		StopGrowth: stop.IfPure,
		SplitScore: impurity.Gini,
		PruneScore: impurity.Gini,
		XValLoss:   loss.Deviance,
		K:          10,
	}
	forest1 := &forest.RandomForestFactory{
		RandomTree: forest.RandomTree{
			NumSamples: 0.7,
			NumLeaves:  20,
			NumColumns: 0,
			SplitScore: impurity.Gini,
		},
		NumTrees: 1000,
	}
	forest2 := &forest.RandomForestFactory{
		RandomTree: forest.RandomTree{
			NumSamples: 0.7,
			NumLeaves:  40,
			NumColumns: 0,
			SplitScore: impurity.Gini,
		},
		NumTrees: 1000,
	}
	methods := []classification.Factory{
		tree1,
		tree2,
		bagging.New(tree1, 4, 0),
		bagging.New(tree1, 16, 0),
		forest1.New(),
		forest2.New(),
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
				name:   sample.GetName(),
				values: make([]<-chan *classification.Result, len(methods)),
			}
			for i, method := range methods {
				r.values[i] = XAssess(method, sample, loss.ZeroOne)
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
		fmt.Println(string([]byte{'A' + byte(i)}), "=", method.GetName())
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
		var errors []string
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

var queue chan int

func XAssess(cf classification.Factory, samples data.Set, L loss.Function) <-chan *classification.Result {
	worker := <-queue
	resChan := make(chan *classification.Result, 1)
	go func() {
		res := classification.Assess(cf, samples, L)
		resChan <- res
		queue <- worker
	}()
	return resChan
}

func init() {
	n := runtime.GOMAXPROCS(0)
	queue = make(chan int, n)
	for i := 0; i < n; i++ {
		queue <- i
	}
}
