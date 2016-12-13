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
	"encoding/json"
	"fmt"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/pkg/errors"
)

var (
	// ErrSegmentNotInExperiment occurs when a user is hashed into a segment that has not been claimed by an experiment.
	ErrSegmentNotInExperiment = errors.New("Segment is not assigned to an experiment")
)

// Namespace is a container for experiments. Segments in the namespace divide
// traffic. Units are the keys that will hash experiments.
type Namespace struct {
	Name        string
	Segments    segments
	Labels      []string
	Experiments []Experiment
}

// NewNamespace creates a new namespace with all segments available. It returns
// an error if no units are given.
func NewNamespace(name string, labels []string) *Namespace {
	n := &Namespace{
		Name:   name,
		Labels: labels,
	}
	return n
}

// ToNamespace is a helper function that converts a Namespace into a
// *storage.Namespace.
func (n *Namespace) ToNamespace() *storage.Namespace {
	ns := &storage.Namespace{
		Name:        n.Name,
		Labels:      n.Labels,
		Experiments: make([]*storage.Experiment, len(n.Experiments)),
	}
	for i, e := range n.Experiments {
		ns.Experiments[i] = e.ToExperiment()
	}
	return ns
}

// AddExperiment adds an experiment to the namespace. It takes the the given
// number of segments from the namespace. It returns an error if the number of
// segments is larger than the number of available segments in the namespace.
// TODO: make sure this is still needed
func (n *Namespace) AddExperiment(e Experiment) error {
	seg, err := n.Segments.Claim(e.Segments)
	if err != nil {
		return errors.Wrap(err, "could not claim segments from namespace")
	}
	n.Segments = seg
	n.Experiments = append(n.Experiments, e)
	return nil
}

// eval evaluates the current namespace given the hashConfig supplied. If the
// user hashes into a segment that is not claimed by an experiment then eval
// returns the error ErrSegmentNotInExperiment. Otherwise it should return an
// ExperimentResponse.
func (n *Namespace) eval(h hashConfig) (ExperimentResponse, error) {
	h.setNs(n.Name)
	i, err := hash(h)
	if err != nil {
		return ExperimentResponse{}, err
	}
	segment := uniform(i, 0, float64(len(n.Segments)*8))
	if !n.Segments.isClaimed(uint64(segment)) {
		return ExperimentResponse{}, ErrSegmentNotInExperiment
	}

	for _, exp := range n.Experiments {
		if !exp.Segments.isClaimed(uint64(segment)) {
			continue
		}
		p, err := exp.eval(h)
		if err != nil {
			return ExperimentResponse{}, err
		}
		return ExperimentResponse{Name: exp.Name, Params: p}, nil

	}

	// unreachable
	return ExperimentResponse{}, fmt.Errorf("unreachable code Namespace.eval")
}

// MarshalJSON implements the json.Marshaler interface for Namespaces.
func (n *Namespace) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Name        string       `json:"name"`
		Segments    segments     `json:"segments"`
		Labels      []string     `json:"labels"`
		Experiments []Experiment `json:"experiments"`
	}{
		Name:        n.Name,
		Segments:    n.Segments,
		Labels:      n.Labels,
		Experiments: n.Experiments,
	}
	return json.Marshal(aux)
}

// ExperimentResponse holds the data for an evaluated expeiment.
type ExperimentResponse struct {
	Name   string
	Params []ParamValue
}

// Namespaces determines the assignments for the a given user's id based on the
// current set of namespaces and experiments. It returns a []ExperimentResponse
// if it is successful or an error if something went wrong.
func (ec *Config) Namespaces(teamID, userID string) ([]ExperimentResponse, error) {
	h := hashConfig{}
	h.setUserID(userID)

	var response []ExperimentResponse
	for _, ns := range TeamNamespaces(*ec.Storage, teamID) {
		eResp, err := ns.eval(h)
		if err == ErrSegmentNotInExperiment {
			continue
		}
		if err != nil {
			return nil, err
		}
		response = append(response, eResp)
	}
	return response, nil
}
