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

import (
	"encoding/hex"
	"math/rand"
	"testing"
)

func TestSegmentContains(t *testing.T) {
	tests := map[string]struct {
		main segments
		in   segments
		want bool
	}{
		"empty matches empty":          {main: segments{make([]byte, 2), 16}, in: segments{make([]byte, 2), 16}, want: true},
		"subset":                       {main: segments{[]byte{0xff, 0xff}, 16}, in: segments{[]byte{0xf0, 0xf0}, 16}, want: true},
		"different lengths fail bytes": {main: segments{[]byte{0, 0}, 16}, in: segments{[]byte{0}, 16}, want: false},
		"different lengtsh fail len":   {main: segments{[]byte{0}, 1}, in: segments{[]byte{0}, 2}, want: false},
		"empty ns non-empty exp":       {main: segments{[]byte{0, 0}, 16}, in: segments{[]byte{255, 255}, 16}, want: false},
		"non claimed":                  {main: segments{[]byte{255, 0}, 16}, in: segments{[]byte{0, 255}, 16}, want: false},
		"overlapping":                  {main: segments{[]byte{0x0f, 0xff}, 16}, in: segments{[]byte{0xf0, 0xf0}, 16}, want: false},
	}
	for k, test := range tests {
		got := test.main.contains(&test.in)
		if got != test.want {
			t.Errorf("%s: %v.contains(%v) = %v want %v", k, test.main, test.in, got, test.want)
		}
	}
}

func TestSegmentIsClaimed(t *testing.T) {
	tests := map[string]struct {
		seg  segments
		n    uint64
		want bool
	}{
		"empty contains 0":   {seg: segments{[]byte{0}, 1}, n: 0, want: false},
		"empty contains 1":   {seg: segments{[]byte{0}, 8}, n: 1, want: false},
		"[1] contains 0":     {seg: segments{[]byte{1}, 8}, n: 0, want: true},
		"[1] contanis 7":     {seg: segments{[]byte{1 << 7}, 8}, n: 7, want: true},
		"[7,12] contains 12": {seg: segments{[]byte{0, 1<<7 | 1<<4}, 16}, n: 12, want: true},
		"[7] contains 14":    {seg: segments{[]byte{0, 1 << 7}, 16}, n: 14, want: false},
	}
	for k, test := range tests {
		got := test.seg.isClaimed(test.n)
		if test.want != got {
			t.Errorf("%s: %v.isClaimed(%v) = %v, want %v", k, test.seg, test.n, got, test.want)
		}
	}
}

func TestSegmentsAvailable(t *testing.T) {
	tests := map[string]struct {
		seg  segments
		want []int
	}{
		"no seg":             {seg: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, want: []int{}},
		"all but one":        {seg: segments{[]byte{127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, want: []int{7}},
		"first bit byte 1,2": {seg: segments{[]byte{127, 127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, want: []int{7, 15}},
		"first byte":         {seg: segments{[]byte{0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, want: []int{0, 1, 2, 3, 4, 5, 6, 7}},
		"second byte":        {seg: segments{[]byte{255, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, want: []int{8, 9, 10, 11, 12, 13, 14, 15}},
		"all first bit": {
			seg:  segments{[]byte{254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254}, 128},
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
		want  segments
	}{
		"set one":      {seg: segments{make([]byte, 16), 128}, index: 0, want: segments{[]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 128}},
		"set thirteen": {seg: segments{make([]byte, 16), 128}, index: 13, want: segments{[]byte{0, 1 << 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 128}},
	}

	for k, test := range tests {
		set(test.seg.b, test.index)
		for i := range test.seg.b {
			if test.seg.b[i] != test.want.b[i] {
				t.Errorf("%s: set(%v %v) = %v, want %v", k, hex.EncodeToString(test.seg.b), test.index, test.seg.b, test.want)
			}
		}
	}
}

func TestSegmentClear(t *testing.T) {
	tests := map[string]struct {
		seg   segments
		index int
		want  segments
	}{
		"unset 15": {seg: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, index: 15, want: segments{[]byte{255, 127, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}},
	}
	for k, test := range tests {
		clear(test.seg.b, test.index)
		for i := range test.seg.b {
			if test.seg.b[i] != test.want.b[i] {
				t.Errorf("%s: clear(%v %v) = %v, want %v", k, hex.EncodeToString(test.seg.b), test.index, test.seg.b, test.want)
			}
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
		"sample one":  {seg: segments{}, num: 1, want: segments{[]byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 128}},
		"sample all":  {seg: segments{}, num: 128, want: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}},
	}

	for k, test := range tests {
		rand.Seed(0)
		got := test.seg.sample(test.num)
		for i := range got {
			if got[i] != test.want.b[i] {
				t.Errorf("%s: test.sample() = %v, want %v", k, got, test.want)
			}
		}
	}
}

func TestSegmentsClaim(t *testing.T) {
	tests := map[string]struct {
		seg  segments
		out  segments
		want segments
		err  error
	}{
		"remove none":          {seg: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, out: segments{}, want: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, err: nil},
		"remove all":           {seg: segments{}, out: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, want: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, err: nil},
		"remove all from none": {seg: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, out: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, want: segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128}, err: ErrSegmentUnavailable},
		"remove some":          {seg: segments{}, out: segments{[]byte{255}, 8}, want: segments{[]byte{255}, 8}, err: nil},
		"bad remove some":      {seg: segments{[]byte{255, 0}, 16}, out: segments{[]byte{255, 0}, 16}, want: segments{[]byte{255, 0}, 16}, err: ErrSegmentUnavailable},
	}

	for k, test := range tests {
		seg, err := test.seg.Claim(&test.out)
		for i := range seg {
			if seg[i] != test.want.b[i] {
				t.Errorf("%s: test.Claim(%v) = %v, want %v", k, test.out, seg, test.want)
			}
			if err != test.err {
				t.Errorf("%s: %v.Remove(%v) = %v, want %v", k, test.seg, test.out, err, test.err)
			}
		}
	}
}
