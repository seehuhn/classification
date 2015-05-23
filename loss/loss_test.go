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
