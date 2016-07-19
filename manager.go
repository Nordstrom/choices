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
	namespaceIndexByTeamID map[string][]int
	namespace              []Namespace
	mu                     sync.RWMutex
}

var defaultManager = &manager{
	namespaceIndexByTeamID: make(map[string][]int, 100),
	namespace:              []Namespace{},
}

func (m *manager) nsByID(teamID string) []int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.namespaceIndexByTeamID[teamID]
}

func (m *manager) ns(index int) Namespace {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.namespace[index]
}

func (m *manager) addns(n *Namespace) error {
	m.namespace = append(m.namespace, *n)
	m.namespaceIndexByTeamID = m.recompute()
	return nil
}

func (m *manager) recompute() map[string][]int {
	out := make(map[string][]int, len(m.namespaceIndexByTeamID))
	for i, n := range m.namespace {
		for _, t := range n.TeamID {
			out[t] = append(out[t], i)
		}
	}
	return out
}

func Addns(n *Namespace) error {
	return defaultManager.addns(n)
}
