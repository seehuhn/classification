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

// Digits represents a data set of digitised, hand-written digits.
// The data set is taken from the book "The Elements of Statistical
// Learning" by Hastie, Tibschirani and Friedman.
var Digits Set = &digits{}

func (d *digits) Name() string {
	return "ZIP code digits"
}

func (d *digits) readFile(fname string) (data *Data, err error) {
	X, Y, err := matrix.Plain.Read(fname, zipDataColumns)
	if err != nil {
		msg := fmt.Sprintf("cannot read %s", fname)
		return nil, &Error{d.Name(), msg, err}
	}
	res := &Data{
		NumClasses: 10,
		X:          X,
		Y:          Y.Column(0),
	}
	return res, nil
}

func (d *digits) TrainingData() (data *Data, err error) {
	return d.readFile(trainFile)
}

func (d *digits) TestData() (data *Data, err error) {
	return d.readFile(testFile)
}
