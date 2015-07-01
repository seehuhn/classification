// matrix.go -
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

var ErrSyntax = errors.New("malformed input data")

type TextFormat struct {
	RowSep   byte
	FieldSep byte
}

var CSVFormat = TextFormat{
	RowSep:   '\n',
	FieldSep: ',',
}

var PlainFormat = TextFormat{
	RowSep:   '\n',
	FieldSep: ' ',
}

type tokenizer struct {
	TextFormat
	scanner     *bufio.Scanner
	atEOL       bool
	lineStarted bool
	row         []float64
}

func newTokenizer(r io.Reader, opts TextFormat) *tokenizer {
	res := &tokenizer{
		scanner:    bufio.NewScanner(r),
		TextFormat: opts,
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
	s.row = []float64{}
	for s.scanner.Scan() {
		token := s.scanner.Text()
		if token == string(s.RowSep) {
			return true
		}
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		x, err := strconv.ParseFloat(token, 64)
		if err != nil && err.(*strconv.NumError).Err == strconv.ErrSyntax {
			x = math.NaN()
		}
		s.row = append(s.row, x)
	}
	return len(s.row) > 0
}

func (s *tokenizer) Row() []float64 {
	return s.row
}

func ReadAsText(fname string, opts TextFormat) (*Float64, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var in io.Reader
	if strings.HasSuffix(fname, ".gz") {
		gin, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		defer gin.Close()
		in = gin
	} else {
		in = file
	}

	n := 0
	p := 0
	data := []float64{}
	scanner := newTokenizer(in, opts)
	for scanner.Scan() {
		row := scanner.Row()
		if p == 0 {
			p = len(row)
		} else if len(row) != p {
			return nil, ErrSyntax
		}
		data = append(data, row...)
		n++
	}

	return NewFloat64(n, p, 0, data), nil
}
