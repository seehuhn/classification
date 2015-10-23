package data

import (
	"fmt"
	"github.com/seehuhn/classification/matrix"
)

const (
	trainFile = "data/zip.train.gz"
	testFile  = "data/zip.test.gz"
)

func zipDataColumns(col int) matrix.ColumnType {
	switch col {
	case 0:
		return matrix.RoundToIntColumn
	case 257: // work around spaces at the end of line
		return matrix.IgnoredColumn
	default:
		return matrix.Float64Column
	}
}

type digits struct{}

var Digits Set = &digits{}

func (d *digits) Name() string {
	return "ZIP code digits"
}

func (d *digits) NumClasses() int {
	return 10
}

func (d *digits) TrainingSet() (X *matrix.Float64, Y []int, err error) {
	XTrain, YTrain, err := matrix.Plain.Read(trainFile, zipDataColumns)
	if err != nil {
		msg := fmt.Sprintf("cannot read %s", trainFile)
		return nil, nil, &Error{d.Name(), msg, err}
	}
	return XTrain, YTrain.Column(0), nil
}

func (d *digits) TestSet() (X *matrix.Float64, Y []int, err error) {
	XTest, YTest, err := matrix.Plain.Read(testFile, zipDataColumns)
	if err != nil {
		msg := fmt.Sprintf("cannot read %s", testFile)
		return nil, nil, &Error{d.Name(), msg, err}
	}
	return XTest, YTest.Column(0), nil
}
