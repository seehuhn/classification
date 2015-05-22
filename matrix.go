package classification

type Matrix struct {
	n, p int
	data []float64
}

func NewMatrix(n, p int, data []float64) *Matrix {
	if data == nil {
		data = make([]float64, n*p)
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
