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

	"k8s.io/apimachinery/pkg/labels"

	"github.com/foolusion/elwinprotos/storage"
)

// ParamValue is a key value pair returned from an evalated experiment
// parameter.
type ParamValue struct {
	Name  string
	Value string
}

// Experiment is a structure that represents a single experiment in elwin. It
// can contain multiple parameters. Experiments are evaluated through the call
// to Namespaces.
type Experiment struct {
	Name     string
	Labels   labels.Set
	Params   []Param
	Segments segments
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
func (e *Experiment) SampleSegments(ns *Namespace, num int) *Experiment {
	seg := ns.Segments.sample(num)
	nsSeg, err := ns.Segments.Claim(seg)
	if err != nil {
		panic(err)
	}
	copy(ns.Segments[:], nsSeg[:])
	e.Segments = seg
	return e
}

// ToExperiment is a helper function that converts an Experiment into a
// *storage.Experiment.
func (e *Experiment) ToExperiment() *storage.Experiment {
	exp := &storage.Experiment{
		Name:     e.Name,
		Labels:   e.Labels,
		Params:   make([]*storage.Param, len(e.Params)),
		Segments: make([]byte, 16),
	}
	copy(exp.Segments[:], e.Segments[:])
	for i, p := range e.Params {
		exp.Params[i] = p.ToParam()
	}
	return exp
}

// eval evaluates the experiment based on the given hashConfig. It returns a
// []ParamValue of the evaluated params or an error.
func (e *Experiment) eval(h hashConfig) ([]ParamValue, error) {
	p := make([]ParamValue, len(e.Params))
	h.salt[2] = e.Name
	for i, param := range e.Params {
		par, err := param.eval(h)
		if err != nil {
			return nil, err
		}
		p[i] = par
	}
	return p, nil
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
	Name  string
	Value Value
}

// ToParam is a helper function that converts a Param into a *storage.Param.
func (p *Param) ToParam() *storage.Param {
	param := &storage.Param{
		Name: p.Name,
	}
	switch val := p.Value.(type) {
	case *Uniform:
		param.Value = &storage.Value{
			Choices: val.Choices,
		}
	case *Weighted:
		param.Value = &storage.Value{
			Choices: val.Choices,
			Weights: val.Weights,
		}
	}
	return param
}

// MarshalJSON implements the json.Marshaler interface for Params.
func (p *Param) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Name   string `json:"name"`
		Values Value  `json:"values"`
	}{
		Name:   p.Name,
		Values: p.Value,
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
	val, err := p.Value.Value(i)
	if err != nil {
		return ParamValue{}, err
	}
	return ParamValue{Name: p.Name, Value: val}, nil
}
