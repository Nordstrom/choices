package choices

import "testing"

func TestHash(t *testing.T) {
	tests := []struct {
		in     []func(*hashConfig)
		out    int64
		outErr error
	}{
		{
			in: []func(*hashConfig){
				hashSalt("hello"),
			},
			out:    769918044520341130,
			outErr: nil,
		},
		{
			in: []func(*hashConfig){
				hashSalt("choices"),
				hashNs("hello"),
				hashExp("test"),
				hashParam("something"),
				hashUnits([]unit{
					{key: "test", value: []string{"one", "two"}},
					{key: "blah", value: []string{"a", "b", "c"}},
				}),
			},
			out:    678521530177676309,
			outErr: nil,
		},
	}

	for _, test := range tests {
		got, err := hash(test.in...)
		if test.out != got || test.outErr != err {
			t.Errorf(
				"hash(%v) = %v, %v, want %v, %v",
				test.in,
				got,
				err,
				test.out,
				test.outErr)
		}
	}
}

func TestUniform(t *testing.T) {
	tests := []struct {
		hash     int64
		min, max float64
		want     float64
	}{
		{
			hash: 0,
			min:  0,
			max:  100,
			want: 0,
		},
		{
			hash: 0xfffffffffffffff,
			min:  0,
			max:  10,
			want: 10,
		},
	}
	for _, test := range tests {
		got := uniform(test.hash, test.min, test.max)
		if test.want != got {
			t.Errorf("uniform(%v, %v, %v) = %v, want %v",
				test.hash,
				test.min,
				test.max,
				got,
				test.want)
		}
	}
}
