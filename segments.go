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

type segments [16]byte

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
