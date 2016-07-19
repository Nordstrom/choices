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
)

var config = struct {
	globalSalt string
}{
	globalSalt: "choices",
}

type Namespace struct {
	Name        string
	Segments    segments
	TeamID      []string
	Experiments []Experiment
	Units       []string
}

func NewNamespace(name, teamID string, units []string) (*Namespace, error) {
	if len(units) == 0 {
		return nil, fmt.Errorf("addns: no units given")
	}
	n := &Namespace{
		Name:     name,
		TeamID:   []string{teamID},
		Units:    units,
		Segments: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
	}
	return n, nil
}

func (n *Namespace) eval(units []unit) ([]paramValue, error) {
	h := &hashConfig{}
	h.hashSalt("choice")
	h.hashNs(n.Name)
	h.hashUnits(units)
	i, err := hash(h)
	if err != nil {
		return nil, err
	}
	segment := uniform(i, 0, float64(len(n.Segments)*8))
	if n.Segments.contains(uint64(segment)) {
		return nil, fmt.Errorf("segment not assigned to an experiment")
	}

	for _, exp := range n.Experiments {
		if !exp.Segments.contains(uint64(segment)) {
			continue
		}
		return exp.eval(h)

	}
	return nil, nil
}

func (n *Namespace) Addexp(name string, params []Param, numSegments int) {
	e := Experiment{
		Name:     name,
		Params:   params,
		Segments: n.Segments.sample(numSegments),
	}
	n.Experiments = append(n.Experiments, e)
}

func filterUnits(units map[string][]string, keep []string) []unit {
	out := make([]unit, len(keep))

	for i, k := range keep {
		out[i] = unit{key: k, value: units[k]}
	}
	return out
}

type Response struct {
	Experiments map[string][]paramValue
}

func (r *Response) add(key string, p []paramValue) {
	if r.Experiments == nil {
		r.Experiments = make(map[string][]paramValue, 10)
	}
	r.Experiments[key] = p
}

func Namespaces(ctx context.Context, m *manager, teamID string, units map[string][]string) (*Response, error) {
	// ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	// defer cancel()

	if m == nil {
		m = defaultManager
	}

	response := &Response{}

	for _, index := range m.nsByID(teamID) {
		ns := m.ns(index)

		u := filterUnits(units, ns.Units)
		r, err := ns.eval(u)
		if err != nil {
			return nil, err
		}
		response.add(ns.Name, r)
	}
	return response, nil
}
