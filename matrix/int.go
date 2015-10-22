// int.go -
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
	"fmt"
	"os"
	"strings"
)

// Int represents matrices of categorical (int) values.  New Int
// matrices can be allocated using the `NewInt` function.
type Int struct {
	n, p, stride int
	data         []int
}

// NewInt allocates a new matrix for categorical (int) values.
// `n` is the number of rows and `p` is the number of columns.  Date
// is stored so that values within one row are consecutive; if
// `stride` is non-zero, each row is followed by `stride-p` unused
// values in memory.  If `data` is non-nil, it is used to store the
// matrix values, values in `data` are preserved and form the initial
// contents of the matrix; otherwise a new slice is allocated for
// storage and all matrix elements are initially `0`.
func NewInt(n, p, stride int, data []int) *Int {
	if data == nil {
		data = make([]int, n*p)
	} else if len(data) < n*p {
		panic("not enough data provided")
	}
	if stride == 0 {
		stride = p
	}
	return &Int{
		n:      n,
		p:      p,
		stride: stride,
		data:   data,
	}
}

// Shape returns the number of rows and column of `mat`.
func (mat *Int) Shape() (int, int) {
	return mat.n, mat.p
}

// At returns the matrix element at row `i`, column `j`.
func (mat *Int) At(i, j int) int {
	return mat.data[i*mat.stride+j]
}

// Set changes the matrix alement at row `i`, column `j` to be `x`.
func (mat *Int) Set(i, j int, x int) {
	mat.data[i*mat.stride+j] = x
}

// Inc increments the matrix alement at row `i`, column `j` by 1.
func (mat *Int) Inc(i, j int) {
	mat.data[i*mat.stride+j]++
}

// Row returns a slice representing row `i` of the matrix.  The
// returned slice is a sub-slice of the matrix data, and any changes
// to elements of the returned row slice are visible in the underlying
// matrix, too.
func (mat *Int) Row(i int) []int {
	base := i * mat.stride
	return mat.data[base : base+mat.p]
}

// Column returns column `j` of the data.  If `mat` has stride 1, the
// returned slice is a sub-slice of the matrix data, and any changes
// to elements of the returned row slice are visible in the underlying
// matrix, too.  Otherwise, the returned slice is a copy of the matrix
// data, and can be changed without changing the original matrix.
func (mat *Int) Column(j int) []int {
	if mat.stride == 1 {
		return mat.data[:mat.n]
	}
	res := make([]int, mat.n)
	for i := 0; i < mat.n; i++ {
		res[i] = mat.At(i, j)
	}
	return res
}

// Format returns a textual, human-readable representation of the
// matrix.  `format` is the format string used for each matrix
// element.
func (mat *Int) Format(format string) string {
	entries := [][]string{}
	for i := 0; i < mat.n; i++ {
		row := []string{}
		for j := 0; j < mat.p; j++ {
			sep := ", "
			if j == mat.p-1 {
				sep = ";"
				if i == mat.n-1 {
					sep = " ]"
				}
			}
			entry := fmt.Sprintf(format, mat.data[i*mat.stride+j]) + sep
			row = append(row, entry)
		}
		entries = append(entries, row)
	}

	for j := 0; j < mat.p; j++ {
		width := 0
		for i := 0; i < mat.n; i++ {
			l := len(entries[i][j])
			if l > width {
				width = l
			}
		}
		for i := 0; i < mat.n; i++ {
			entries[i][j] += strings.Repeat(" ", width-len(entries[i][j]))
		}
	}

	rows := []string{}
	for i, rowEntries := range entries {
		head := "  "
		if i == 0 {
			head = "[ "
		}
		rows = append(rows, head+strings.Join(rowEntries, ""))
	}
	return strings.Join(rows, "\n")
}

// String returns a textual, human-readable representation of the
// matrix.
func (mat *Int) String() string {
	return mat.Format("%6d")
}

// WriteCSV writes the matrix in .csv form into the file with name
// `fname`.  Any pre-existing file with this name is over-written.
func (mat *Int) WriteCSV(fname string) {
	fd, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for i := 0; i < mat.n; i++ {
		for j := 0; j < mat.p; j++ {
			if j > 0 {
				fd.WriteString(",")
			}
			fmt.Fprintf(fd, "%d", mat.At(i, j))
		}
		fd.WriteString("\n")
	}
}
