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

import "math/rand"

type segments [16]byte

// SegmentsAll is a value where every segment is available
var SegmentsAll = segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}

func (s *segments) contains(seg uint64) bool {
	index, pos := seg/8, seg%8
	return s[index]>>pos&1 == 1
}

func (s *segments) available() []int {
	out := make([]int, 0, 128)
	for i := range s {
		for shift := uint8(0); shift < 8; shift++ {
			if s[i]&(1<<shift) == 1<<shift {
				out = append(out, i*8+int(shift))
			}
		}
	}
	return out
}

type bit int

const (
	zero bit = iota
	one
)

func (s *segments) set(index int, val bit) {
	i, pos := index/8, uint8(index%8)
	switch val {
	case zero:
		s[i] &= ^(1 << pos)
	case one:
		s[i] |= 1 << pos
	}
}

func (s *segments) sample(n int) segments {
	avail := s.available()
	out := segments{}
	p := rand.Perm(len(avail))
	for i := 0; i < n; i++ {
		s.set(avail[p[i]], zero)
		out.set(avail[p[i]], one)
	}
	return out
}

func (s *segments) count() int {
	count := 0
	for _, v := range s {
		count += int(cnt[v])
	}
	return count
}

var cnt = [256]byte{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8}
