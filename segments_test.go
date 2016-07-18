package choices

import "testing"

func TestSegmentContains(t *testing.T) {
	tests := []struct {
		seg  segments
		n    uint64
		want bool
	}{
		{
			seg:  segments{},
			n:    0,
			want: false,
		},
		{
			seg:  segments{},
			n:    1,
			want: false,
		},
		{
			seg:  segments{1},
			n:    0,
			want: true,
		},
		{
			seg:  segments{1 << 7},
			n:    7,
			want: true,
		},
		{
			seg:  segments{0, 1<<7 | 1<<4},
			n:    12,
			want: true,
		},
		{
			seg:  segments{0, 1 << 7},
			n:    14,
			want: false,
		},
	}
	for _, test := range tests {
		got := test.seg.contains(test.n)
		if test.want != got {
			t.Errorf("%v.contains(%v) = %v, want %v",
				test.seg,
				test.n,
				got,
				test.want)
		}
	}
}

func TestSegmentsAvailable(t *testing.T) {
	tests := []struct {
		seg  segments
		want []int
	}{
		{seg: segments{}, want: []int{}},
		{seg: segments{1}, want: []int{0}},
		{seg: segments{1, 1}, want: []int{0, 8}},
		{seg: segments{255, 1}, want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8}},
		{seg: segments{1, 255}, want: []int{0, 8, 9, 10, 11, 12, 13, 14, 15}},
		{seg: segments{1, 127}, want: []int{0, 8, 9, 10, 11, 12, 13, 14}},
		{
			seg:  segments{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			want: []int{0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120},
		},
		{
			seg:  segments{1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7, 1 << 7},
			want: []int{7, 15, 23, 31, 39, 47, 55, 63, 71, 79, 87, 95, 103, 111, 119, 127},
		},
	}
	for _, test := range tests {
		got := test.seg.available()
		if len(got) != len(test.want) {
			t.Errorf("%v.available() = %v, want %v", test.seg, got, test.want)
			t.FailNow()
		}
		for i, v := range got {
			if v != test.want[i] {
				t.Errorf("%v.available() = %v, want %v", test.seg, got, test.want)
				t.FailNow()
			}
		}

	}
}
