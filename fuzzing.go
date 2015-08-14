// +build gofuzz

package classification

import (
	"bytes"
)

func Fuzz(data []byte) int {
	r := bytes.NewReader(data)
	tree, err := TreeFromFile(r)
	if err != nil {
		if tree != nil {
			panic("tree is not nil")
		}
		return 0
	}

	w := &bytes.Buffer{}
	err = tree.WriteTo(w)
	if err != nil {
		panic(err)
	}
	data2 := w.Bytes()

	r = bytes.NewReader(data2)
	tree2, err := TreeFromFile(r)
	if err != nil {
		panic(err)
	}
	w = &bytes.Buffer{}
	err = tree2.WriteTo(w)
	if err != nil {
		panic(err)
	}
	data3 := w.Bytes()

	if bytes.Compare(data2, data3) != 0 {
		panic("re-encoded tree differs from original")
	}

	return 1
}
