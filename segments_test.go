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

import (
	"math/rand"
	"testing"
)

func TestSegmentContains(t *testing.T) {
	tests := map[string]struct {
		seg  segments
		n    uint64
		want bool
	}{
		"empty contains 0":   {seg: segments{}, n: 0, want: false},
		"empty contains 1":   {seg: segments{}, n: 1, want: false},
		"[1] contains 0":     {seg: segments{1}, n: 0, want: true},
		"[1] contanis 7":     {seg: segments{1 << 7}, n: 7, want: true},
		"[7,12] contains 12": {seg: segments{0, 1<<7 | 1<<4}, n: 12, want: true},
		"[7] contains 17":    {seg: segments{0, 1 << 7}, n: 14, want: false},
	}
	for k, test := range tests {
		got := test.seg.claimed(test.n)
		if test.want != got {
			t.Errorf("%s: %v.contains(%v) = %v, want %v", k, test.seg, test.n, got, test.want)
		}
	}
}

func TestSegmentsAvailable(t *testing.T) {
	tests := map[string]struct {
		seg  segments
		want []int
	}{
		"no seg":             {seg: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: []int{}},
		"all but one":        {seg: segments{127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: []int{7}},
		"first bit byte 1,2": {seg: segments{127, 127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: []int{7, 15}},
		"first byte":         {seg: segments{0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: []int{0, 1, 2, 3, 4, 5, 6, 7}},
		"second byte":        {seg: segments{255, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: []int{8, 9, 10, 11, 12, 13, 14, 15}},
		"all first bit": {
			seg:  segments{254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254},
			want: []int{0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120},
		},
	}
	for key, test := range tests {
		got := test.seg.available()
		if len(got) != len(test.want) {
			t.Errorf("%s: %v.available() = %v, want %v", key, test.seg, got, test.want)
			t.FailNow()
		}
		for i, v := range got {
			if v != test.want[i] {
				t.Errorf("%s: %v.available() = %v, want %v", key, test.seg, got, test.want)
				t.FailNow()
			}
		}

	}
}

func TestSegmentsSet(t *testing.T) {
	tests := map[string]struct {
		seg   segments
		index int
		value bit
		want  segments
	}{
		"set one":      {seg: segments{}, index: 0, value: one, want: segments{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		"set thirteen": {seg: segments{}, index: 13, value: one, want: segments{0, 1 << 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		"unset 15":     {seg: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, index: 15, value: zero, want: segments{255, 127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}},
	}

	for k, test := range tests {
		test.seg = test.seg.set(test.index, test.value)
		if test.seg != test.want {
			t.Errorf("%s: test.set(%v, %v) = %v, want %v", k, test.index, test.value, test.seg, test.want)
		}
	}
}

func TestSegmentsSample(t *testing.T) {
	tests := map[string]struct {
		seg  segments
		num  int
		want segments
	}{
		"sample none": {seg: segments{}, num: 0, want: segments{}},
		"sample one":  {seg: segments{}, num: 1, want: segments{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}

	for k, test := range tests {
		rand.Seed(0)
		_, got := test.seg.sample(test.num)
		if got != test.want {
			t.Errorf("%s: test.sample() = %v, want %v", k, got, test.want)
		}
	}
}

func TestSegmentsRemove(t *testing.T) {
	tests := map[string]struct {
		seg  segments
		out  segments
		want segments
		err  error
	}{
		"remove none":          {seg: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, out: segments{}, want: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, err: nil},
		"remove all":           {seg: segments{}, out: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, err: nil},
		"remove all from none": {seg: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, out: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, err: ErrSegmentUnavailable},
		"remove some":          {seg: segments{}, out: segments{255}, want: segments{255}, err: nil},
		"bad remove some":      {seg: segments{255, 0}, out: segments{255, 0}, want: segments{255, 0}, err: ErrSegmentUnavailable},
	}

	for k, test := range tests {
		seg, err := test.seg.Claim(test.out)
		if seg != test.want {
			t.Errorf("%s: test.Claim(%v) = %v, want %v", k, test.out, seg, test.want)
		}
		if err != test.err {
			t.Errorf("%s: %v.Remove(%v) = %v, want %v", k, test.seg, test.out, err, test.err)
		}
	}
}
