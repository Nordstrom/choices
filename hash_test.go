// Copyright 2016 Andrew O'Neill, Nordstrom

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
		in     hashConfig
		out    uint64
		outErr error
	}{
		{
			in: hashConfig{
				salt: [3]string{"", "", ""},
			},
			out:    11350825000285427928,
			outErr: nil,
		},
		{
			in: hashConfig{
				salt:   [3]string{"hello", "test", "something"},
				userID: "my-user-id",
			},
			out:    13108963186855807482,
			outErr: nil,
		},
	}

	for _, test := range tests {
		got, err := hash(test.in)
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
		{
			hash: 0x8000000000000000,
			min:  0,
			max:  100,
			want: 50,
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
	backup := globalSalt
	defer func() {
		globalSalt = backup
	}()
	globalSalt = "salt"
	h := hashConfig{
		salt:   [3]string{"namespace", "experiment", "param"},
		userID: "abcdef1234567890",
	}
	for i := 0; i < b.N; i++ {
		if _, err := hash(h); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashBytes(b *testing.B) {
	h := hashConfig{}
	for i := 0; i < b.N; i++ {
		if _, err := h.Bytes(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashBytesAll(b *testing.B) {
	backup := globalSalt
	defer func() {
		globalSalt = backup
	}()
	globalSalt = "salt"
	h := hashConfig{salt: [3]string{"namespace", "experiment", "param"}, userID: "value"}
	for i := 0; i < b.N; i++ {
		if _, err := h.Bytes(); err != nil {
			b.Fatal(err)
		}
	}
}
