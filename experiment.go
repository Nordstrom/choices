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
	"errors"

	"github.com/foolusion/elwinprotos/storage"
	"github.com/pquerna/ffjson/ffjson"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	// ErrSegmentNotInExperiment occurs when a user is hashed into a
	// segment that has not been claimed by an experiment.
	ErrSegmentNotInExperiment = errors.New("Segment is not assigned to an experiment")
)

// Experiment is a structure that represents a single experiment in
// elwin. It can contain multiple parameters.
type Experiment struct {
	ID        string
	Name      string
	Namespace string
	Labels    labels.Set
	Params    []Param
	Segments  *segments
}

// NewExperiment creates an experiment with the supplied name and no
// segments claimed. In order for any traffic to be assigned to this
// experiment you will need to call Experiment.SetSegments or
// Experiment.SampleSegments.
func NewExperiment(name string) *Experiment {
	return &Experiment{
		Name: name,
	}
}

// ToExperiment is a helper function that converts an Experiment into
// a *storage.Experiment. It will overwrite the provided
// *storage.Experiment or create a new one if it is nil
func (e *Experiment) ToExperiment(exp *storage.Experiment) *storage.Experiment {
	if exp == nil {
		exp = new(storage.Experiment)
	}
	exp.Id = e.ID
	exp.Name = e.Name
	exp.Namespace = e.Namespace
	exp.Labels = e.Labels
	exp.Segments = e.Segments.ToSegments(exp.Segments)

	if exp.Params == nil || len(e.Params) != len(exp.Params) {
		exp.Params = make([]*storage.Param, len(e.Params))
	}
	for i, p := range e.Params {
		exp.Params[i] = p.ToParam(exp.Params[i])
	}
	return exp
}

// eval evaluates the experiment based on the given hashConfig. It
// returns a []ParamValue of the evaluated params or an error.
func (e *Experiment) eval(er *ExperimentResponse, h hashConfig) error {
	h.salt[1] = e.Namespace
	i, err := hash(h)
	if err != nil {
		return err
	}
	segment := uniform(i, 0, float64(e.Segments.len))
	if !e.Segments.isClaimed(uint64(segment)) {
		return ErrSegmentNotInExperiment
	}

	if er == nil {
		er = new(ExperimentResponse)
		er.Params = make([]*ParamValue, len(e.Params))
	} else if er.Params == nil {
		er.Params = make([]*ParamValue, len(e.Params))
	}
	h.salt[2] = e.Name
	for i, param := range e.Params {
		if er.Params[i] == nil {
			er.Params[i] = new(ParamValue)
		}
		err := param.eval(er.Params[i], h)
		if err != nil {
			return err
		}
	}
	if len(er.Params) > len(e.Params) {
		er.Params = er.Params[:len(e.Params)]
	}
	er.Name = e.Name
	er.Namespace = e.Namespace
	er.Labels = e.Labels
	return nil
}

// MarshalJSON implements the json.Marshaler interface for
// Experiments.
func (e *Experiment) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Name     string    `json:"name"`
		Segments *segments `json:"segments"`
		Params   []Param   `json:"params"`
	}{
		Name:     e.Name,
		Segments: e.Segments,
		Params:   e.Params,
	}
	return ffjson.Marshal(aux)
}

// Param is a struct that represents a single parameter in an
// experiment. Param is evaluated through the call to Experiments.
type Param struct {
	Name    string
	Choices choice
}

// ToParam is a helper function that converts a Param into a
// *storage.Param.
func (p *Param) ToParam(param *storage.Param) *storage.Param {
	if param == nil {
		param = new(storage.Param)
	}
	param.Name = p.Name

	if param.Value == nil {
		param.Value = new(storage.Value)
	}

	switch val := p.Choices.(type) {
	case *Uniform:
		param.Value.Choices = val.Choices
	case *Weighted:
		c := make([]string, len(val.Choices))
		weights := make([]float64, len(val.Choices))
		for i, v := range val.Choices {
			c[i] = v.name
			weights[i] = v.weight
		}
		param.Value.Choices = c
		param.Value.Weights = weights
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
	return ffjson.Marshal(aux)
}

// eval evaluates the Param based on the given hashConfig. It returns
// a ParamValue containing the value the user is assigned.
func (p *Param) eval(ep *ParamValue, h hashConfig) error {
	if ep == nil {
		return errors.New("param value is nil")
	}
	h.setParam(p.Name)
	i, err := hash(h)
	if err != nil {
		return err
	}
	val, err := p.Choices.Choice(i)
	if err != nil {
		return err
	}
	*ep = ParamValue{Name: p.Name, Value: val}
	return nil
}

// ExperimentResponse holds the data for an evaluated experiment.
type ExperimentResponse struct {
	Name      string
	Namespace string
	Params    []*ParamValue
	Labels    map[string]string
}

// ParamValue is a key value pair returned from an evaluated experiment
// parameter.
type ParamValue struct {
	Name  string
	Value string
}

// Experiments determines the assignments for the a given user's id
// based on the current set of namespaces and experiments. It returns
// a []*ExperimentResponse if it is successful or an error if something
// went wrong.
func (c *Config) Experiments(out []*ExperimentResponse, userID string, selector labels.Selector) ([]*ExperimentResponse, error) {
	var h hashConfig
	h.setUserID(userID)
	experiments := c.storage.read()
	if out == nil {
		out = make([]*ExperimentResponse, len(experiments))
	} else {
		out = out[:len(experiments)]
	}

	erindex := 0
	for i := range experiments {
		if !selector.Matches(experiments[i].Labels) {
			continue
		}
		if out[erindex] == nil {
			out[erindex] = new(ExperimentResponse)
		}
		err := experiments[i].eval(out[erindex], h)
		if err == ErrSegmentNotInExperiment {
			continue
		}
		if err != nil {
			return nil, err
		}
		erindex++
	}
	return out, nil
}
