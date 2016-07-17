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
	Name       string
	Definition Definition
	Segments   segments
}

func (e *Experiment) eval(ns string, units []unit) ([]paramValue, error) {
	p := make([]paramValue, len(e.Definition.Params))
	for _, param := range e.Definition.Params {
		p = append(p, param.eval(ns, e.Name, units))
	}
	return p, nil
}

type Definition struct {
	Params []Param
}

type Param struct {
	Name  string
	Type  string
	Value Value
}

func (p *Param) eval(ns, exp string, units []unit) paramValue {
	return paramValue{name: p.Name, value: p.Value.Value(ns, exp, p.Name, units)}
}
