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

type unit struct {
	key   string
	value []string
}

type hashConfig struct {
	salt       string
	namespace  string
	experiment string
	param      string
	units      []unit
}

func (h *hashConfig) hashSalt(s string) {
	h.salt = s
}

func (h *hashConfig) hashNs(ns string) {
	h.namespace = ns
}

func (h *hashConfig) hashExp(exp string) {
	h.experiment = exp
}

func (h *hashConfig) hashParam(p string) {
	h.param = p
}

func (h *hashConfig) hashUnits(u []unit) {
	h.units = u
}

func (h *hashConfig) Bytes() []byte {
	var buf bytes.Buffer

	addString := func(s string) {
		if s != "" {
			if buf.Len() != 0 {
				buf.WriteByte('.')
			}
			buf.WriteString(s)
		}
	}

	addString(h.salt)
	addString(h.namespace)
	addString(h.experiment)
	addString(h.param)
	if len(h.units) != 0 {
		if buf.Len() != 0 {
			buf.WriteByte('.')
		}
		for _, unit := range h.units {
			buf.WriteString(unit.key)
			buf.WriteByte('=')
			for i, v := range unit.value {
				buf.WriteString(v)
				if i < len(unit.value)-1 {
					buf.WriteByte(',')
				}
			}
		}
	}
	return buf.Bytes()
}

func hash(h *hashConfig) (uint64, error) {
	if h == nil {
		h = &hashConfig{}
	}
	hash := sha1.Sum(h.Bytes())
	i := binary.BigEndian.Uint64(hash[:8])
	return i, nil
}

func uniform(hash uint64, min, max float64) float64 {
	return min + (max-min)*(float64(hash)/longScale)
}
