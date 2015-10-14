// matrix.go - matrices for continuous and categorical data
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

package matrix

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// ErrSyntax is returned by the `TextFormat.Read` method when the
// input file is malformed.
var ErrSyntax = errors.New("malformed input data")

// TextFormat describes the file format for a textual representation
// of a matrix.
type TextFormat struct {
	// RowSep specifies the character which separates matrix rows in
	// the input.  Normally this will be '\n', so that matrix rows
	// correspond to lines in the input file.
	RowSep byte

	// FieldSep specifies the character which separates matrix
	// elements within a row.
	FieldSep byte
}

// CSV describes the file format for .csv files (comma-separated
// values).
var CSV = &TextFormat{
	RowSep:   '\n',
	FieldSep: ',',
}

// Plain describes the file format where matrix entries within a row
// are separated by spaces, and matrix rows correspond to rows in the
// text file.
var Plain = &TextFormat{
	RowSep:   '\n',
	FieldSep: ' ',
}

func (opts *TextFormat) Read(fname string, cols ColumnFunc) (*Float64, *Int, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var in io.Reader
	if strings.HasSuffix(fname, ".gz") {
		gin, err := gzip.NewReader(file)
		if err != nil {
			return nil, nil, err
		}
		defer gin.Close()
		in = gin
	} else {
		in = file
	}

	n := 0
	pFloat64 := 0
	pInt := 0
	float64Data := []float64{}
	intData := []int{}
	scanner := newTokenizer(in, cols, opts)
	for scanner.Scan() {
		float64Row, intRow := scanner.Row()
		if n == 0 {
			pFloat64 = len(float64Row)
			pInt = len(intRow)
		} else if len(float64Row) != pFloat64 || len(intRow) != pInt {
			return nil, nil, ErrSyntax
		}
		float64Data = append(float64Data, float64Row...)
		intData = append(intData, intRow...)
		n++
	}

	return NewFloat64(n, pFloat64, 0, float64Data), NewInt(n, pInt, 0, intData), nil
}

// ColumnType is used by `matrix.ColumnFunc` to specify the role of
// individual columns in the input file.
type ColumnType int

const (
	// Float64Column indicates columns for continuous inputs,
	// represented by `float64` values in the program.
	Float64Column ColumnType = iota

	// IntColumn indicates columns for categorical inputs, represented
	// by `int` values in the program.
	IntColumn

	// IgnoredColumn indicates columns in the input file which should
	// be ignored.
	IgnoredColumn
)

// ColumnFunc is the type of functions used to determine the type of
// each column.  The argument of the ColumnFunc is the column index,
// starting with `0` for the first column.
type ColumnFunc func(int) ColumnType

type tokenizer struct {
	*TextFormat
	scanner     *bufio.Scanner
	atEOL       bool
	lineStarted bool

	cols       ColumnFunc
	float64Row []float64
	intRow     []int
}

func newTokenizer(r io.Reader, cols ColumnFunc, opts *TextFormat) *tokenizer {
	res := &tokenizer{
		scanner:    bufio.NewScanner(r),
		TextFormat: opts,
		cols:       cols,
	}
	res.scanner.Split(res.splitField)
	return res
}

func (s *tokenizer) splitField(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if s.atEOL {
		s.atEOL = false
		s.lineStarted = false
		return 0, []byte{s.RowSep}, nil
	}
	start := 0
	for i, c := range data {
		if i == start && c == s.RowSep && !s.lineStarted {
			start = i + 1
			continue
		}
		if c == s.RowSep {
			s.atEOL = true
			s.lineStarted = true
			return i + 1, data[start:i], nil
		} else if c == s.FieldSep {
			s.lineStarted = true
			return i + 1, data[start:i], nil
		}
	}
	data = data[start:]
	if atEOF {
		if s.lineStarted {
			s.atEOL = true
		}
		if len(data) == 0 {
			return 0, nil, io.EOF
		}
		return len(data), data, nil
	}
	return
}

func (s *tokenizer) Scan() bool {
	s.float64Row = []float64{}
	s.intRow = []int{}
	col := 0
	for s.scanner.Scan() {
		token := s.scanner.Text()
		if token == string(s.RowSep) {
			return true
		}
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		colType := s.cols(col)
		col++
		switch colType {
		case Float64Column:
			x, err := strconv.ParseFloat(token, 64)
			if err != nil && err.(*strconv.NumError).Err == strconv.ErrSyntax {
				x = math.NaN()
			}
			s.float64Row = append(s.float64Row, x)
		case IntColumn:
			x, err := strconv.ParseFloat(token, 64)
			if err != nil {
				// TODO(voss): how to return an error to the caller here?
				panic(err)
			}
			s.intRow = append(s.intRow, int(x))
		}
	}
	return len(s.intRow) > 0 || len(s.float64Row) > 0
}

func (s *tokenizer) Row() ([]float64, []int) {
	return s.float64Row, s.intRow
}
