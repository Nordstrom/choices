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

type paramValue struct {
	name  string
	value string
}

type Experiment struct {
	Name     string
	Params   []Param
	Segments segments
}

func (e *Experiment) eval(h hashConfig) ([]paramValue, error) {
	p := make([]paramValue, 0, len(e.Params))
	h.setExp(e.Name)
	for _, param := range e.Params {
		par, err := param.eval(h)
		if err != nil {
			return nil, err
		}
		p = append(p, par)
	}
	return p, nil
}

type Param struct {
	Name  string
	Value Value
}

func (p *Param) eval(h hashConfig) (paramValue, error) {
	h.setParam(p.Name)
	val, err := p.Value.Value(h)
	if err != nil {
		return paramValue{}, err
	}
	return paramValue{name: p.Name, value: val}, nil
}
