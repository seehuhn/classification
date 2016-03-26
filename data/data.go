// data.go - a data structure to represent a data set
//
// Copyright (C) 2015  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package data

import (
	"github.com/seehuhn/classification/matrix"
)

// Data represents a collection of samples with associated classes,
// used either as trainings data or as test data in a classification
// problem.
type Data struct {
	// NumClasses gives the number of classes used in the data set.
	NumClasses int

	// X stores the predictor variables.  Each column of the matrix
	// corresponds to a variable, rows represent samples.
	X *matrix.Float64

	// Y stores the response variables.  The length of Y must equal
	// the number of rows of X.
	Y []int

	// Weights is either nil or a vector of the same length as Y.  If
	// the vector is present, it must the the same length as Y and
	// stores weights for each sample.  If `Weight` is nil, each
	// sample has weight 1.x
	Weights []float64

	// Rows can be either nil or a vector or row numbers.  If `Rows`
	// is non-nil, the entries must be valid row numbers for `X` and
	// the data set consists of the listed rows.  If `Rows` is nil,
	// the data set consists of all rows of `X`.
	Rows []int
}

// NRow returns the number of samples in the data set.
func (data *Data) NRow() int {
	if data.Rows != nil {
		return len(data.Rows)
	}
	return len(data.Y)
}

// NCol returns the number of predictor variables in the data set.
func (data *Data) NCol() int {
	_, p := data.X.Shape()
	return p
}

// GetRows returns the vector of all rows used in the data set.  The
// returned slice is owned by the Data object and must not be changed
// by the caller.
func (data *Data) GetRows() []int {
	if data.Rows != nil {
		return data.Rows
	}
	res := make([]int, len(data.Y))
	for i := range res {
		res[i] = i
	}
	return res
}
