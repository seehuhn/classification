package matrix

import (
	"bufio"
	"math"
	"strings"
	"testing"
)

func TestReadTokenize(t *testing.T) {
	testData := []struct {
		in  string
		out []string
	}{
		{"a bc def", []string{"a", "bc", "def", "\n"}},
		{"a\nbc def", []string{"a", "\n", "bc", "def", "\n"}},
		{"a bc def\n", []string{"a", "bc", "def", "\n"}},
		{"\na bc def\n\n\n", []string{"a", "bc", "def", "\n"}},
		{"\n", []string{}},
		{"", []string{}},
	}
	for _, data := range testData {
		r := strings.NewReader(data.in)
		scanner := bufio.NewScanner(r)
		tokenizer := &tokenizer{
			TextFormat: TextFormat{
				RowSep:   '\n',
				FieldSep: ' ',
			},
		}
		scanner.Split(tokenizer.splitField)
		pos := 0
		for scanner.Scan() {
			if pos >= len(data.out) {
				t.Errorf("too much output from %q", data.in)
				break
			}
			token := scanner.Text()
			if token != data.out[pos] {
				t.Errorf("wrong token %d for %q: expected %q, got %q",
					pos+1, data.in, data.out[pos], token)
				break
			}
			pos++
		}
		if pos < len(data.out) {
			t.Errorf("missing output from %q, expected %q next",
				data.in, data.out[pos])
		}
	}
}

func floatSlicesEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if (!math.IsNaN(b[i]) || !math.IsNaN(x)) && b[i] != x {
			return false
		}
	}
	return true
}

func TestReadRow(t *testing.T) {
	testData := []struct {
		in  string
		out [][]float64
	}{
		{"1,2,3\n4,5,6\n", [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}},
		{"1,2,3\n", [][]float64{{1.0, 2.0, 3.0}}},
		{"1,2,3", [][]float64{{1.0, 2.0, 3.0}}},
		{"1\n2\n3\n", [][]float64{{1.0}, {2.0}, {3.0}}},
		{"1\n2\n3", [][]float64{{1.0}, {2.0}, {3.0}}},
		{"1,2\n\n", [][]float64{{1.0, 2.0}}},
		{"1, 2\n", [][]float64{{1.0, 2.0}}},
		{"1,x,3\n", [][]float64{{1.0, math.NaN(), 3.0}}},
	}
	for _, data := range testData {
		r := strings.NewReader(data.in)
		scanner := newTokenizer(r,
			func(int) ColumnType { return Float64Column },
			TextFormat{'\n', ','})
		pos := 0
		for scanner.Scan() {
			row, _ := scanner.Row()
			if !floatSlicesEqual(row, data.out[pos]) {
				t.Errorf("wrong scanner output for %q: expected %v, got %v",
					data.in, data.out[pos], row)
			}
			pos++
		}
	}
}
