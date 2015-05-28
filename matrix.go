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

package classification

import (
	"fmt"
	"os"
)

type Matrix struct {
	n, p int
	data []float64
}

func NewMatrix(n, p int, data []float64) *Matrix {
	if data == nil {
		data = make([]float64, n*p)
	} else if len(data) < n*p {
		panic("not enough data provided")
	}
	return &Matrix{
		n:    n,
		p:    p,
		data: data,
	}
}

func (mat *Matrix) At(i, j int) float64 {
	return mat.data[i*mat.p+j]
}

func (mat *Matrix) Row(i int) []float64 {
	return mat.data[i*mat.p : (i+1)*mat.p]
}

func (mat *Matrix) WriteCSV(fname string) {
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
			fmt.Fprintf(fd, "%g", mat.At(i, j))
		}
		fd.WriteString("\n")
	}
}
