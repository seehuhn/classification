// files.go -
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

func WriteVector(fname string, x []int) {
	fd, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for _, xi := range x {
		fmt.Fprintf(fd, "%d\n", xi)
	}
}
