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
	"errors"
	"math/rand"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/pquerna/ffjson/ffjson"
)

type segments struct {
	b   []byte
	len int
}

func (s *segments) ToSegments(seg *storage.Segments) *storage.Segments {
	if seg == nil {
		seg = new(storage.Segments)
	}
	seg.B = s.b
	seg.Len = int64(s.len)

	return seg
}

// segmentsAll is a value where every segment is available
var segmentsAll = segments{b: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, len: 128}

var (
	// ErrSegmentUnavailable is thrown when you request an a segment set to
	// 0, an unavailable segment.
	ErrSegmentUnavailable = errors.New("segment unavailable")
	// ErrDifferentLengthSegments is thrown when you make a request with
	// segments that have different lengths.
	ErrDifferentLengthSegments = errors.New("segments have different lengths")
)

func (s *segments) contains(in *segments) bool {
	if len(s.b) != len(in.b) || s.len != in.len {
		return false
	}
	for i := range s.b {
		if s.b[i]&in.b[i] != in.b[i] {
			return false
		}
	}
	return true
}

// Claim claims the segments in out from s and throws an error if a
// segment has already been claimed. If there is no error you can
// overwrite the segments.b with the returned value.
func (s *segments) Claim(out *segments) ([]byte, error) {
	if len(s.b) != len(out.b) || s.len != out.len {
		return nil, ErrDifferentLengthSegments
	}
	b := make([]byte, len(out.b))
	for i := range b {
		if s.b[i]&out.b[i] > 0 {
			return s.b, ErrSegmentUnavailable
		}
		b[i] = s.b[i] | out.b[i]
	}
	return b, nil
}

// isClaimed returns whether or not a given segment is claimed from the
// segments.
func (s *segments) isClaimed(seg uint64) bool {
	index, pos := seg/8, seg%8
	return s.b[index]>>pos&1 == 1
}

// available returns a list of segments that are available.
func (s *segments) available() []int {
	out := make([]int, 0, s.len)
	for i := range s.b {
		for shift := uint8(0); shift < 8; shift++ {
			// check if we have reached the end
			if i*8+int(shift) > s.len {
				break
			}
			if s.b[i]&(1<<shift) != 1<<shift {
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

func set(b []byte, index int) {
	i, pos := index/8, uint8(index%8)
	b[i] |= 1 << pos
}

func clear(b []byte, index int) {
	i, pos := index/8, uint8(index%8)
	b[i] &= ^(1 << pos)
}

// sample samples segments that are unclaimed and returns a the segments and
// the sample. If n is greater than the number of available segments it returns
// an error.
func (s *segments) sample(n int) []byte {
	avail := s.available()
	p := rand.Perm(len(avail))
	out := make([]byte, len(s.b))
	if n > len(avail) {
		n = len(avail)
	}
	for i := 0; i < n; i++ {
		set(out, avail[p[i]])
	}
	return out
}

// count returns the number of claimed segments
func (s *segments) count() int {
	count := 0
	for _, v := range s.b {
		count += int(cnt[v])
	}
	return count
}

// MarshalJSON implements the json.Marshaler interface for segments
func (s *segments) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(hex.EncodeToString(s.b[:]))
}

var cnt = [256]byte{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8}
