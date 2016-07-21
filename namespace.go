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
	"fmt"

	"github.com/foolusion/choices/elwin"
)

var config = struct {
	globalSalt string
}{
	globalSalt: "choices",
}

// Namespace is a container for experiments. Segments in the namespace divide
// traffic. Units are the keys that will hash experiments.
type Namespace struct {
	Name        string
	Segments    segments
	TeamID      []string
	Experiments []Experiment
	Units       []string
}

// NewNamespace creates a new namespace with all segments available. It returns
// an error if no units are given.
func NewNamespace(name, teamID string, units []string) (*Namespace, error) {
	if len(units) == 0 {
		return nil, fmt.Errorf("addns: no units given")
	}
	n := &Namespace{
		Name:     name,
		TeamID:   []string{teamID},
		Units:    units,
		Segments: segmentsAll,
	}
	return n, nil
}

func (n *Namespace) eval(userID string) (*elwin.Experiment, error) {
	h := hashConfig{}
	h.setSalt(config.globalSalt)
	h.setNs(n.Name)
	h.setUserID(userID)
	i, err := hash(h)
	if err != nil {
		return nil, err
	}
	segment := uniform(i, 0, float64(len(n.Segments)*8))
	if n.Segments.contains(uint64(segment)) {
		return nil, nil
	}

	for _, exp := range n.Experiments {
		if !exp.Segments.contains(uint64(segment)) {
			continue
		}
		return exp.eval(h)

	}
	return nil, fmt.Errorf("shouldn't be here")
}

// Addexp adds an experiment to the namespace. It takes the the given number of
// segments from the namespace. It returns an error if the number of segments
// is larger than the number of available segments in the namespace.
func (n *Namespace) Addexp(name string, params []Param, numSegments int) error {
	if n.Segments.count() < numSegments {
		return fmt.Errorf("Namespace.Addexp: not enough segments in namespace, want: %v, got %v", numSegments, n.Segments.count())
	}
	e := Experiment{
		Name:     name,
		Params:   params,
		Segments: n.Segments.sample(numSegments),
	}
	n.Experiments = append(n.Experiments, e)
	return nil
}

func filterUnits(units map[string][]string, keep []string) []unit {
	out := make([]unit, len(keep))

	for i, k := range keep {
		out[i] = unit{key: k, value: units[k]}
	}
	return out
}

// Response is what is returned to the user on a call to Namespaces.
type Response struct {
	Experiments map[string][]paramValue
}

func (r *Response) add(key string, p []paramValue) {
	if r.Experiments == nil {
		r.Experiments = make(map[string][]paramValue, 10)
	}
	r.Experiments[key] = p
}

// Namespaces determines the assignments for the a given users units based on
// the current set of namespaces and experiments. It returns a Response object
// if it is successful or an error if something went wrong.
func Namespaces(m *manager, teamID, userID string) (*elwin.Experiments, error) {
	if m == nil {
		m = defaultManager
	}

	response := &elwin.Experiments{}

	for _, index := range m.nsByID(teamID) {
		ns := m.ns(index)

		r, err := ns.eval(userID)
		if err != nil {
			return nil, err
		}
		response.Experiments[ns.Name] = r
	}
	return response, nil
}
