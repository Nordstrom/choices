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

func hashSalt(s string) func(*hashConfig) {
	return func(h *hashConfig) {
		if h == nil {
			h = &hashConfig{}
		}
		h.salt = s
	}
}

func hashNs(ns string) func(*hashConfig) {
	return func(h *hashConfig) {
		if h == nil {
			h = &hashConfig{}
		}
		h.namespace = ns
	}
}

func hashExp(exp string) func(*hashConfig) {
	return func(h *hashConfig) {
		if h == nil {
			h = &hashConfig{}
		}
		h.experiment = exp
	}
}

func hashParam(p string) func(*hashConfig) {
	return func(h *hashConfig) {
		if h == nil {
			h = &hashConfig{}
		}
		h.param = p
	}
}

func hashUnits(u []unit) func(*hashConfig) {
	return func(h *hashConfig) {
		if h == nil {
			h = &hashConfig{}
		}
		h.units = u
	}
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

func hash(funcs ...func(*hashConfig)) (uint64, error) {
	h := &hashConfig{}
	for _, f := range funcs {
		f(h)
	}

	hash := sha1.Sum(h.Bytes())
	i := binary.BigEndian.Uint64(hash[:8])
	return i, nil
}

func uniform(hash uint64, min, max float64) float64 {
	return min + (max-min)*(float64(hash)/longScale)
}
