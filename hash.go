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
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
)

const longScale = float64(0xFFFFFFFFFFFFFFFF)

var (
	// globalSalt should only ever be set once.
	globalSalt = "choices"
)

var (
	calledSet bool

	// ErrGlobalSaltAlreadySet is the error returned when a
	// SetGlobalSalt has been called more than once.
	ErrGlobalSaltAlreadySet = errors.New("global salt already set")
)

// SetGlobalSalt this sets the global salt used for hashing users. It
// should only ever be called once. It returns an error if it is
// called more than once.
func SetGlobalSalt(s string) error {
	if calledSet {
		return ErrGlobalSaltAlreadySet
	}

	globalSalt = s
	calledSet = true
	return nil
}

// WithGlobalSalt is a configuration option for Config that sets the
// globalSalt to something other than the default.
func WithGlobalSalt(s string) ConfigOpt {
	return func(c *Config) error {
		globalSalt = s
		return nil
	}
}

// hashConfig is a struct to store the hash data.
type hashConfig struct {
	// salt has 3 parts namespace is index 0, experiment is index 1, and param is index 2
	salt   [3]string
	userID string
}

// HashExperience takes the supplied arguments and returns the hashed
// uint64 that can be used for determining a segment.
func HashExperience(namespace, experiment, param, userID string) (uint64, error) {
	h := hashConfig{userID: userID}
	h.setNs(namespace)
	h.setExp(experiment)
	h.setParam(param)
	return hash(h)
}

// setNs sets the namespaces portion of the salt
func (h *hashConfig) setNs(ns string) {
	h.salt[0] = ns
}

// setExp sets the experiment portion of the salt
func (h *hashConfig) setExp(exp string) {
	h.salt[1] = exp
}

// setParam sets the param portion of the salt
func (h *hashConfig) setParam(p string) {
	h.salt[2] = p
}

// setUserID sets the userID
func (h *hashConfig) setUserID(u string) {
	h.userID = u
}

// errWriter is a convenience struct to eliminate lots of if err !=
// nil.
type errWriter struct {
	buf bytes.Buffer
	err error
}

// writeString writes the given string to the buf or nothing if err !=
// nil
func (e *errWriter) writeString(s string) {
	if e.err == nil {
		_, e.err = e.buf.WriteString(s)
	}
}

// writeByte writes a single byte to the buf or nothing if err != nil
func (e *errWriter) writeByte(b byte) {
	if e.err == nil {
		e.err = e.buf.WriteByte(b)
	}
}

// Bytes returns a []byte that represents the entire salt+userID the
// format is as follows.
//     "globalSalt.namespace.experiment.param@userID"
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

// hash hashes the hashConfig returning a uint64 hashed value or an
// error.
func hash(h hashConfig) (uint64, error) {
	b, err := h.Bytes()
	if err != nil {
		return 0, err
	}
	hash := sha1.Sum(b)
	i := binary.BigEndian.Uint64(hash[:8])
	return i, nil
}

// uniform returns a uniformly random value between the min and max
// values supplied.
func uniform(hash uint64, min, max float64) float64 {
	return min + (max-min)*(float64(hash)/longScale)
}
