// loss_test.go -
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

package loss

import (
	"testing"
)

var all = []Function{
	Deviance,
	ZeroOne,
	Other,
}

func TestLoss(t *testing.T) {
	for _, fn := range all {
		for p := 1; p < 10; p++ {
			if p > 1 {
				prob := make([]float64, p)
				for i := range prob {
					prob[i] = 1 / float64(p)
				}
				for i := range prob {
					loss := fn(i, prob)
					if loss <= 0 {
						t.Error("unexpected loss value", loss, "in the interior")
					}
				}
			}

			for i := 0; i < p; i++ {
				prob := make([]float64, p)
				prob[i] = 1
				loss := fn(i, prob)
				if loss < 0 || loss > 1e-6 {
					t.Error("unexpected loss value", loss, "at the boundary")
				}
			}
		}
	}
}
