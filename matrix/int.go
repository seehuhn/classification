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

type Int struct {
	n, p, stride int
	data         []int
}

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

func (mat *Int) At(i, j int) int {
	return mat.data[i*mat.stride+j]
}

func (mat *Int) Row(i int) []int {
	base := i * mat.stride
	return mat.data[base : base+mat.p]
}

func (mat *Int) Column(j int) []int {
	if mat.stride == 1 {
		return mat.data
	}
	res := make([]int, mat.n)
	for i := 0; i < mat.n; i++ {
		res[i] = mat.At(i, j)
	}
	return res
}

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

func (mat *Int) String() string {
	return mat.Format("%6d")
}

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
