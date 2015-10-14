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
	"math"
	"testing"
)

var all = []Function{
	Deviance,
	ZeroOne,
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

func TestZeroOne(t *testing.T) {
	hist := []float64{1, 3, 2}
	l0 := ZeroOne(0, hist)
	l1 := ZeroOne(1, hist)
	l2 := ZeroOne(2, hist)
	if math.Abs(l0-1.0) > 1e-6 || math.Abs(l1-0.0) > 1e-6 || math.Abs(l2-1.0) > 1e-6 {
		t.Error("unexpected loss values, expected 1 0 1, got", l0, l1, l2)
	}

	hist = []float64{1, 2, 2}
	l0 = ZeroOne(0, hist)
	l1 = ZeroOne(1, hist)
	l2 = ZeroOne(2, hist)
	if math.Abs(l0-1.0) > 1e-6 || math.Abs(l1-0.5) > 1e-6 || math.Abs(l2-0.5) > 1e-6 {
		t.Error("unexpected loss values, expected 1 0.5 0.5, got", l0, l1, l2)
	}

	hist = []float64{1, 1, 1}
	l0 = ZeroOne(0, hist)
	l1 = ZeroOne(1, hist)
	l2 = ZeroOne(2, hist)
	q := 2.0 / 3.0
	if math.Abs(l0-q) > 1e-6 || math.Abs(l1-q) > 1e-6 || math.Abs(l2-q) > 1e-6 {
		t.Error("unexpected loss values, expected 1/3 1/3 1/3, got", l0, l1, l2)
	}
}
