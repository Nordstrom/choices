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
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
)

type segments [16]byte

// segmentsAll is a value where every segment is available
var segmentsAll = segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}

var (
	// ErrSegmentUnavailable is thrown when you request an a segment set to
	// 0, an unavailable segment.
	ErrSegmentUnavailable = errors.New("segment unavailable")
)

// Claim removes the segments in del from s and throws an error if the
func (s segments) Claim(out segments) (segments, error) {
	var seg segments
	for i := range seg {
		if s[i]&out[i] > 0 {
			return s, ErrSegmentUnavailable
		}
		seg[i] = s[i] | out[i]
	}
	return seg, nil
}

// isClaimed returns whether or not a given segment is claimed from the
// segments.
func (s segments) isClaimed(seg uint64) bool {
	index, pos := seg/8, seg%8
	return s[index]>>pos&1 == 1
}

// available returns a list of segments that are available.
func (s segments) available() []int {
	out := make([]int, 0, 128)
	for i := range s {
		for shift := uint8(0); shift < 8; shift++ {
			if s[i]&(1<<shift) != 1<<shift {
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

// set claims or releases the segment at the given index
func (s segments) set(index int, val bit) segments {
	i, pos := index/8, uint8(index%8)
	switch val {
	case zero:
		s[i] &= ^(1 << pos)
	case one:
		s[i] |= 1 << pos
	}
	return s
}

// sample samples segments that are unclaimed and returns a the segments and
// the sample. If n is greater than the number of available segments it returns
// an error.
func (s segments) sample(n int) segments {
	avail := s.available()
	p := rand.Perm(len(avail))
	var out segments
	if n > len(avail) {
		n = len(avail)
	}
	for i := 0; i < n; i++ {
		out = out.set(avail[p[i]], one)
	}
	return out
}

// count returns the number of claimed segments
func (s segments) count() int {
	count := 0
	for _, v := range s {
		count += int(cnt[v])
	}
	return count
}

// InSegment takes a namespace name, userID and segments and returns whether
// the userID is in a segment that is claimed. Typical use case for this is
// determining whether or not a user is in a given experiment.
func InSegment(namespace, userID string, s segments) bool {
	h := hashConfig{userID: userID}
	h.setNs(namespace)
	i, err := hash(h)
	if err != nil {
		return false
	}
	segment := uniform(i, 0, float64(len(s)*8))

	return s.isClaimed(uint64(segment))
}

// MarshalJSON implements the json.Marshaler interface for segments
func (s segments) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(s[:]))
}

var cnt = [256]byte{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8}
