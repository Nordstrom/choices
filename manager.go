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

import "sync"

type manager struct {
	namespace []Namespace
	mu        sync.RWMutex
}

var defaultStorage = &manager{
	namespace: []Namespace{},
}

func (m *manager) TeamNamespaces(teamID string) []Namespace {
	ns := make([]Namespace, 0, len(m.namespace))
	m.mu.RLock()
	for _, n := range m.namespace {
		for _, t := range n.TeamID {
			if t == teamID {
				ns = append(ns, n)
			}
		}
	}
	m.mu.RUnlock()
	return ns
}

func (m *manager) addns(n *Namespace) error {
	m.mu.Lock()
	m.namespace = append(m.namespace, *n)
	m.mu.Unlock()
	return nil
}

func Addns(n *Namespace) error {
	return defaultStorage.addns(n)
}
