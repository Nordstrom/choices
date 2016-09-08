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
	"bytes"
	"crypto/sha1"
	"encoding/binary"
)

const longScale = float64(0xFFFFFFFFFFFFFFFF)

type hashConfig struct {
	salt   [4]string
	userID string
}

func (h *hashConfig) setSalt(s string) {
	h.salt[0] = s
}

func (h *hashConfig) setNs(ns string) {
	h.salt[1] = ns
}

func (h *hashConfig) setExp(exp string) {
	h.salt[2] = exp
}

func (h *hashConfig) setParam(p string) {
	h.salt[3] = p
}

func (h *hashConfig) setUserID(u string) {
	h.userID = u
}

type errWriter struct {
	buf bytes.Buffer
	err error
}

func (e *errWriter) writeString(s string) {
	if e.err == nil {
		_, e.err = e.buf.WriteString(s)
	}
}

func (e *errWriter) writeByte(b byte) {
	if e.err == nil {
		e.err = e.buf.WriteByte(b)
	}
}

func (h *hashConfig) Bytes() ([]byte, error) {
	ew := errWriter{}

	for i, v := range h.salt {
		ew.writeString(v)
		if i < len(h.salt)-1 {
			ew.writeByte('.')
		}
	}

	ew.writeByte('@')

	ew.writeString(h.userID)

	if ew.err != nil {
		return nil, ew.err
	}
	return ew.buf.Bytes(), nil
}

func hash(h hashConfig) (uint64, error) {
	b, err := h.Bytes()
	if err != nil {
		return 0, err
	}
	hash := sha1.Sum(b)
	i := binary.BigEndian.Uint64(hash[:8])
	return i, nil
}

func uniform(hash uint64, min, max float64) float64 {
	return min + (max-min)*(float64(hash)/longScale)
}
