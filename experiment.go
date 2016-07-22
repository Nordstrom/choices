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

import "github.com/foolusion/choices/elwin"

type paramValue struct {
	Experiment string
	Name       string
	Value      string
}

type Experiment struct {
	Name     string
	Params   []Param
	Segments segments
}

func (e *Experiment) eval(h hashConfig) (*elwin.Experiment, error) {
	p := make([]*elwin.Param, len(e.Params))
	h.salt[2] = e.Name
	for i, param := range e.Params {
		par, err := param.eval(h)
		if err != nil {
			return nil, err
		}
		p[i] = par
	}
	return &elwin.Experiment{Params: p}, nil
}

type Param struct {
	Name  string
	Value Value
}

func (p *Param) eval(h hashConfig) (*elwin.Param, error) {
	h.salt[3] = p.Name
	i, err := hash(h)
	if err != nil {
		return nil, err
	}
	val, err := p.Value.Value(i)
	if err != nil {
		return nil, err
	}
	return &elwin.Param{Key: p.Name, Value: val}, nil
}
