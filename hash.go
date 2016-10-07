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
	"errors"
)

const longScale = float64(0xFFFFFFFFFFFFFFFF)

var (
	globalSalt = "choices"
)

var (
	calledSet               bool
	ErrGlobalSaltAlreadySet = errors.New("global salt already set")
)

func SetGlobalSalt(s string) error {
	if calledSet {
		return ErrGlobalSaltAlreadySet
	}

	globalSalt = s
	calledSet = true
	return nil
}

func WithGlobalSalt(s string) ConfigOpt {
	return func(c *Config) error {
		globalSalt = s
		return nil
	}
}

type hashConfig struct {
	salt   [3]string
	userID string
}

func HashExperience(namespace, experiment, param, userID string) (uint64, error) {
	h := hashConfig{userID: userID}
	h.setNs(namespace)
	h.setExp(experiment)
	h.setParam(param)
	return hash(h)
}

func (h *hashConfig) setNs(ns string) {
	h.salt[0] = ns
}

func (h *hashConfig) setExp(exp string) {
	h.salt[1] = exp
}

func (h *hashConfig) setParam(p string) {
	h.salt[2] = p
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
	ew.writeString(globalSalt)
	ew.writeByte('.')
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
