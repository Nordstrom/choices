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

func WithMemStore(m *MemStore) func(*choices.Config) error {
	return func(e *choices.Config) error {
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

// ExampleData creates a MemStore with some example data loaded into it.
func ExampleData() *MemStore {
	m := &MemStore{}
	t1 := choices.NewNamespace("t1", "test")
	t1.AddExperiment(
		"uniform",
		[]choices.Param{{Name: "a", Value: &choices.Uniform{Choices: []string{"b", "c"}}}},
		128,
	)
	m.addns(t1)
	t2 := choices.NewNamespace("t2", "test")
	t2.AddExperiment(
		"weighted",
		[]choices.Param{{Name: "b", Value: &choices.Weighted{Choices: []string{"on", "off"}, Weights: []float64{2, 1}}}},
		128,
	)
	m.addns(t2)
	t3 := choices.NewNamespace("t3", "test")
	t3.AddExperiment(
		"halfSegments",
		[]choices.Param{{Name: "b", Value: &choices.Uniform{Choices: []string{"on"}}}},
		64,
	)
	m.addns(t3)
	t4 := choices.NewNamespace("t4", "test")
	t4.AddExperiment(
		"multi",
		[]choices.Param{
			{Name: "a", Value: &choices.Uniform{Choices: []string{"on", "off"}}},
			{Name: "b", Value: &choices.Weighted{Choices: []string{"up", "down"}, Weights: []float64{1, 2}}},
		},
		128,
	)
	m.addns(t4)
	return m
}
