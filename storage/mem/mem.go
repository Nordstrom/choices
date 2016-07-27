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

package mem

import (
	"sync"

	"github.com/foolusion/choices"
)

func WithMemStore(m *MemStore) func(*choices.ElwinConfig) error {
	return func(e *choices.ElwinConfig) error {
		e.Storage = m
		return nil
	}
}

type MemStore struct {
	namespace []choices.Namespace
	nsmu      sync.RWMutex

	next    []choices.Namespace
	changed bool
	nextmu  sync.RWMutex
}

func (m *MemStore) Read() []choices.Namespace {
	m.nsmu.RLock()
	ns := m.namespace
	m.nsmu.RUnlock()
	return ns
}

func (m *MemStore) Update() error {
	m.nextmu.RLock()
	if !m.changed {
		m.nextmu.RUnlock()
		return nil
	}
	a := make([]choices.Namespace, len(m.next))
	copy(a, m.next)
	m.nextmu.RUnlock()

	m.nsmu.Lock()
	m.nextmu.Lock()
	m.namespace = a
	m.changed = false
	m.nextmu.Unlock()
	m.nsmu.Unlock()
	return nil
}

func (m *MemStore) addns(n *choices.Namespace) {
	m.nextmu.Lock()
	m.next = append(m.next, *n)
	m.changed = true
	m.nextmu.Unlock()
}

// AddNamespace adds the given namespace to the defaultStorage implementation.
// This is just a basic list of namespaces that is concurrency safe.
func (m *MemStore) AddNamespace(n *choices.Namespace) {
	m.addns(n)
}
