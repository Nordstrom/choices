package choices

import (
	"fmt"
	"testing"
)

func TestValueUniform(t *testing.T) {
}

func TestValueWeighted(t *testing.T) {
	w := &Weighted{Choices: []string{"a", "b", "c"}, Weights: []float64{1, 2, 1}}
	tests := []struct {
		in   uint64
		want string
	}{
		{
			in:   0,
			want: "a",
		},
		{
			in:   0x8000000000000000,
			want: "b",
		},
		{
			in:   0xffffffffffffffff,
			want: "c",
		},
	}

	for _, test := range tests {
		got, _ := w.Value(test.in)
		if got != test.want {
			fmt.Errorf("%v.Value(%v) = %v, want %v", *w, test.in, got, test.want)
		}
	}
}
