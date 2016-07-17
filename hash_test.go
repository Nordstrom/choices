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
