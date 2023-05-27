package main

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"

	"seehuhn.de/go/classification"
	"seehuhn.de/go/classification/bagging"
	"seehuhn.de/go/classification/data"
	"seehuhn.de/go/classification/forest"
	"seehuhn.de/go/classification/impurity"
	"seehuhn.de/go/classification/loss"
	"seehuhn.de/go/classification/tree"
	"seehuhn.de/go/classification/tree/stop"
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
	fmt.Println(" training[s] test[s]")
	fmt.Print(strings.Repeat("-", sampleNameLength))
	for range methods {
		fmt.Print("+" + strings.Repeat("-", colWidth-1))
	}
	fmt.Println(" --------------------")

	methodTrainingTime := make([]time.Duration, len(methods))
	methodTestTime := make([]time.Duration, len(methods))
	for row := range rows {
		var rowTrainingTime, rowTestTime time.Duration
		fmt.Printf("%-*s", sampleNameLength, row.name)
		var errors []string
		for i, c := range row.values {
			value := <-c
			if value.Err != nil {
				fmt.Print("| ERROR" + strings.Repeat(" ", colWidth-7))
				errors = append(errors, value.Err.Error())
			} else {
				fmt.Print("| " + formatVal(value.MeanLoss, value.StdErr,
					colWidth-2, maxPrec))
			}
			rowTrainingTime += value.TrainingTime
			methodTrainingTime[i] += value.TrainingTime
			rowTestTime += value.TestTime
			methodTestTime[i] += value.TestTime
		}
		fmt.Printf("      %6.1f  %6.1f\n",
			rowTrainingTime.Seconds(), rowTestTime.Seconds())
		for _, msg := range errors {
			fmt.Println("  " + msg)
		}
	}
	fmt.Println()
	for i := range methods {
		fmt.Printf("%s = %6.1f  %6.1f\n", string([]byte{'A' + byte(i)}),
			methodTrainingTime[i].Seconds(), methodTestTime[i].Seconds())
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
