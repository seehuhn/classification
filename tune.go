// tune.go -
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
	"math"
)

type tuneProfile struct {
	breaks []float64
	values []float64
}

func (p *tuneProfile) Get(pos float64) float64 {
	i := 0
	n := len(p.breaks)
	for i < n && p.breaks[i] < pos {
		i++
	}
	return p.values[i]
}

func (p *tuneProfile) Minimum() (pos, val float64) {
	bestVal := math.Inf(+1)
	bestI := -1
	for i, val := range p.values {
		if val <= bestVal {
			bestI = i
			bestVal = val
		}
	}
	if bestI <= 0 {
		return 0.0, bestVal
	} else if bestI >= len(p.breaks) {
		return math.Inf(+1), bestVal
	} else {
		return (p.breaks[bestI-1] + p.breaks[bestI]) / 2.0, bestVal
	}
}

func (p *tuneProfile) Add(breaks []float64, values []float64) {
	ni := len(p.breaks)
	nj := len(breaks)
	newBreaks := make([]float64, ni+nj)
	newValues := make([]float64, ni+nj+1)
	if len(p.values) == 0 {
		copy(newBreaks, breaks)
		copy(newValues, values)
	} else {
		newValues[0] = p.values[0] + values[0]
		i := 0
		j := 0
		pos := 0
		for i < ni || j < nj {
			takeI := (j == nj) || (i < ni && p.breaks[i] < breaks[j])
			if takeI {
				newBreaks[pos] = p.breaks[i]
				i++
			} else {
				newBreaks[pos] = breaks[j]
				j++
			}
			pos++
			newValues[pos] = p.values[i] + values[j]
		}
	}
	p.breaks = newBreaks
	p.values = newValues
}
