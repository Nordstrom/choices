// Copyright 2016 Andrew O'Neill

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package choices

import "testing"

func TestHash(t *testing.T) {
	tests := []struct {
		in     []func(*hashConfig)
		out    uint64
		outErr error
	}{
		{
			in: []func(*hashConfig){
				hashSalt("hello"),
			},
			out:    12318688712325458082,
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
			out:    10856344482842820951,
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
		hash     uint64
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
			hash: 0xffffffffffffffff,
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

func BenchmarkHash(b *testing.B) {
	funcs := []func(*hashConfig){
		hashSalt("salt"),
		hashNs("namespace"),
		hashExp("experiment"),
		hashParam("param"),
		hashUnits([]unit{
			{key: "userid", value: []string{"abcdef1234567890"}},
		}),
	}
	for i := 0; i < b.N; i++ {
		hash(funcs...)
	}
}

func BenchmarkHashBytes(b *testing.B) {
	h := hashConfig{}
	for i := 0; i < b.N; i++ {
		h.Bytes()
	}
}

func BenchmarkHashBytesAll(b *testing.B) {
	h := hashConfig{salt: "salt", namespace: "namespace", experiment: "experiment", param: "param", units: []unit{{key: "key", value: []string{"value"}}}}
	for i := 0; i < b.N; i++ {
		h.Bytes()
	}
}
