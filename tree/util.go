package tree

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
