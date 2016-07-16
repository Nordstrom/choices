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
	"context"
	"fmt"
	"sync"
	"time"
)

var config = struct {
	globalSalt string
}{
	globalSalt: "choices",
}

var expManager = struct {
	namespaceIndexByTeamID map[string][]int
	namespace              []Namespace
	mu                     sync.RWMutex
}{
	namespaceIndexByTeamID: make(map[string][]int, 100),
	namespace:              []Namespace{},
}

func nsByID(teamID string) []int {
	expManager.mu.RLock()
	defer expManager.mu.RUnlock()

	return expManager.namespaceIndexByTeamID[teamID]
}

func ns(index int) Namespace {
	expManager.mu.RLock()
	defer expManager.mu.RUnlock()

	return expManager.namespace[index]
}

type Namespace struct {
	Name        string
	Segments    segments
	TeamID      []string
	Experiments []Experiment
	Units       []string
}

func (n *Namespace) eval(units []unit) (int, error) {
	// TODO: determine the segment
	i, err := hash(hashNs(n.Name), hashUnits(units))
	if err != nil {
		return 0, err
	}
	segment := uniform(i, 0, float64(len(n.Segments)*8))
	// TODO: if segment in available segment return
	if n.Segments.contains(uint64(segment)) {
		return 0, fmt.Errorf("segment not assigned to an experiment")
	}

	// TODO: for each experiment
	for _, exp := range n.Experiments {
		// TODO:   if segment not in experiment continue
		if !exp.Segments.contains(uint64(segment)) {
			continue
		}
		// TODO:   return eval experiment
		return exp.eval(n.Name, units)

	}
	// TODO: return default value
	return 0, nil
}

func filterUnits(units map[string][]string, keep []string) []unit {
	out := make([]unit, len(keep))

	for i, k := range keep {
		out[i] = unit{key: k, value: units[k]}
	}
	return out
}

type Response struct{}

func (r *Response) add(i int) {}

func Namespaces(ctx context.Context, teamID string, units map[string][]string) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	response := &Response{}

	namespaces := nsByID(teamID)
	for _, index := range namespaces {
		ns := ns(index)

		u := filterUnits(units, ns.Units)
		r, err := ns.eval(u)
		if err != nil {
			return nil, err
		}
		response.add(r)
	}
	return response, nil
}

type Experiment struct {
	Name       string
	Definition Definition
	Segments   segments
}

func (e *Experiment) eval(ns string, units []unit) (int, error) {
	// TODO: implement this
	return 0, nil
}

type Definition struct {
	Params []Param
}

type Param struct {
	Name  string
	Type  string
	Value Value
}

type Value interface {
	Value(context.Context) string
	String() string
}

type Uniform struct {
	Choices []string
	choice  int
}

func (u *Uniform) Value() string {
	u.eval()
	return u.Choices[u.choice]
}

func (u *Uniform) String() string {
	return u.Choices[u.choice]
}

func (u *Uniform) eval() error {
	i, err := hash()
	if err != nil {
		return err
	}
	u.choice = int(i) % len(u.Choices)
	return nil
}

type Weighted struct {
	Choices []string
	Weights []float64
	choice  int
}

func (w *Weighted) Value() string {
	w.eval()
	return w.Choices[w.choice]
}

func (w *Weighted) String() string {
	return w.Choices[w.choice]
}

func (w *Weighted) eval() error {
	if len(w.Choices) != len(w.Weights) {
		return fmt.Errorf("len(w.Choices) != len(w.Weights): %v != %v", len(w.Choices), len(w.Weights))
	}

	i, err := hash()
	if err != nil {
		return err
	}

	selection := make([]float64, len(w.Weights))
	cumSum := 0.0
	for i, v := range w.Weights {
		cumSum += v
		selection[i] = cumSum
	}
	choice := uniform(i, 0, cumSum)
	for i, v := range selection {
		if choice < v {
			w.choice = i
			return nil
		}
	}

	return fmt.Errorf("no choice made")
}
