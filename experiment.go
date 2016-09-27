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

import storage "github.com/foolusion/choices/elwinstorage"

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
	Params   []Param
	Segments segments
}

func (e *Experiment) ToExperiment() *storage.Experiment {
	exp := &storage.Experiment{
		Name:   e.Name,
		Params: make([]*storage.Param, len(e.Params)),
	}
	copy(exp.Segments[:], e.Segments[:])
	for i, p := range e.Params {
		exp.Params[i] = p.ToParam()
	}
	return exp
}

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

// Param is a struct that represents a single parameter in an experiment. Param
// is evaluated through the call to Namespaces.
type Param struct {
	Name  string
	Value Value
}

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

func (p *Param) eval(h hashConfig) (ParamValue, error) {
	h.salt[3] = p.Name
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
