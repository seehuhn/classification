package classification

func intRange(n int) []int {
	res := make([]int, n)
	for i := range res {
		res[i] = i
	}
	return res
}

func copyIntSlice(src []int) []int {
	res := make([]int, len(src))
	copy(res, src)
	return res
}

func applyRows(data []int, rows []int) []int {
	res := make([]int, len(rows))
	for i, row := range rows {
		res[i] = data[row]
	}
	return res
}
