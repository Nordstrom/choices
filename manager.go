package choices

import (
	"fmt"
	"sync"
)

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

func (m *manager) addns(name, teamID string, units []string) error {
	if len(units) == 0 {
		return fmt.Errorf("addns: no units given")
	}
	i := len(m.namespace)
	n := Namespace{
		Name:     name,
		TeamID:   []string{teamID},
		Units:    units,
		Segments: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
	}
	m.namespace = append(m.namespace, n)
	m.namespaceIndexByTeamID[teamID] = append(m.namespaceIndexByTeamID[teamID], i)
	return nil
}
