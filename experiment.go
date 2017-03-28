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
	"errors"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/foolusion/elwinprotos/storage"
)

var (
	// ErrSegmentNotInExperiment occurs when a user is hashed into a
	// segment that has not been claimed by an experiment.
	ErrSegmentNotInExperiment = errors.New("Segment is not assigned to an experiment")
)

// Experiment is a structure that represents a single experiment in elwin. It
// can contain multiple parameters. Experiments are evaluated through the call
// to Namespaces.
type Experiment struct {
	id        string
	Name      string
	Namespace string
	Labels    labels.Set
	Params    []Param
	Segments  segments
}

// NewExperiment creates an experiment with the supplied name and no segments
// cliamed. In order for any traffic to be assigned to this experiment you will
// need to call Experiment.SetSegments or Experiment.SampleSegments.
func NewExperiment(name string) *Experiment {
	return &Experiment{
		Name: name,
	}
}

// SetSegments copies the segments supplied to the experiment.
func (e *Experiment) SetSegments(seg segments) *Experiment {
	copy(e.Segments[:], seg[:])
	return e
}

// SampleSegments takes a namespace and an amount of segments you want in your
// experiment and returns a random sample of the unclaimed segments from the
// namespace.
func (e *Experiment) SampleSegments(ns segments, num int) segments {
	seg := ns.sample(num)
	nsSeg, err := ns.Claim(seg)
	if err != nil {
		panic(err)
	}
	if i := copy(ns[:], nsSeg[:]); i != 16 {
		panic("didn't do full copy")
	}
	return seg
}

// ToExperiment is a helper function that converts an Experiment into a
// *storage.Experiment.
func (e *Experiment) ToExperiment() *storage.Experiment {
	exp := &storage.Experiment{
		Id:        e.id,
		Name:      e.Name,
		Namespace: e.Namespace,
		Labels:    e.Labels,
		Params:    make([]*storage.Param, len(e.Params)),
		Segments:  make([]byte, len(e.Segments)),
	}
	copy(exp.Segments[:], e.Segments[:])
	for i, p := range e.Params {
		exp.Params[i] = p.ToParam()
	}
	return exp
}

// eval evaluates the experiment based on the given hashConfig. It returns a
// []ParamValue of the evaluated params or an error.
func (e *Experiment) eval(h hashConfig) (ExperimentResponse, error) {
	h.salt[1] = e.Namespace
	i, err := hash(h)
	if err != nil {
		return ExperimentResponse{}, err
	}
	segment := uniform(i, 0, float64(len(e.Segments)*8))
	if !e.Segments.isClaimed(uint64(segment)) {
		return ExperimentResponse{}, ErrSegmentNotInExperiment
	}

	p := make([]ParamValue, len(e.Params))
	h.salt[2] = e.Name
	for i, param := range e.Params {
		par, err := param.eval(h)
		if err != nil {
			return ExperimentResponse{}, err
		}
		p[i] = par
	}
	return ExperimentResponse{
		Name:      e.Name,
		Namespace: e.Namespace,
		Params:    p,
	}, nil
}

// MarshalJSON implements the json.Marshaler interface for Experiments.
func (e *Experiment) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Name     string   `json:"name"`
		Segments segments `json:"segments"`
		Params   []Param  `json:"params"`
	}{
		Name:     e.Name,
		Segments: e.Segments,
		Params:   e.Params,
	}
	return json.Marshal(aux)
}

// Param is a struct that represents a single parameter in an experiment. Param
// is evaluated through the call to Namespaces.
type Param struct {
	Name    string
	Choices choice
}

// ToParam is a helper function that converts a Param into a *storage.Param.
func (p *Param) ToParam() *storage.Param {
	param := &storage.Param{
		Name: p.Name,
	}
	switch val := p.Choices.(type) {
	case *Uniform:
		param.Value = &storage.Value{
			Choices: val.Choices,
		}
	case *Weighted:
		choices := make([]string, len(val.Choices))
		weights := make([]float64, len(val.Choices))
		for i, v := range val.Choices {
			choices[i] = v.name
			weights[i] = v.weight
		}
		param.Value = &storage.Value{
			Choices: choices,
			Weights: weights,
		}
	}
	return param
}

// MarshalJSON implements the json.Marshaler interface for Params.
func (p *Param) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Name    string `json:"name"`
		Choices choice `json:"choices"`
	}{
		Name:    p.Name,
		Choices: p.Choices,
	}
	return json.Marshal(aux)
}

// eval evaluates the Param based on the given hashConfid. It returns a
// ParamValue containing the value the user is assigned.
func (p *Param) eval(h hashConfig) (ParamValue, error) {
	h.setParam(p.Name)
	i, err := hash(h)
	if err != nil {
		return ParamValue{}, err
	}
	val, err := p.Choices.Choice(i)
	if err != nil {
		return ParamValue{}, err
	}
	return ParamValue{Name: p.Name, Value: val}, nil
}

// ExperimentResponse holds the data for an evaluated expeiment.
type ExperimentResponse struct {
	Name      string
	Namespace string
	Params    []ParamValue
}

// ParamValue is a key value pair returned from an evalated experiment
// parameter.
type ParamValue struct {
	Name  string
	Value string
}

// Namespaces determines the assignments for the a given user's id based on the
// current set of namespaces and experiments. It returns a []ExperimentResponse
// if it is successful or an error if something went wrong.
func (ec *Config) Namespaces(userID string, selector labels.Selector) ([]ExperimentResponse, error) {
	h := hashConfig{}
	h.setUserID(userID)

	var response []ExperimentResponse
	for _, exp := range teamNamespaces(ec.storage, selector) {
		eResp, err := exp.eval(h)
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
