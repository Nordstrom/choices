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
	"fmt"

	"github.com/pquerna/ffjson/ffjson"
)

// ValueType are the different types of Values a Param can have. This is used
// for parsing params in storage.
type ValueType int

// Constants for ValueTypes
const (
	ValueTypeBad ValueType = iota
	ValueTypeUniform
	ValueTypeWeighted
)

// choice is the interface Param Choices must implement. They take a hash value
// and return the string that represents the value or an error.
type choice interface {
	Choice(i uint64) (string, error)
}

// Uniform is a way to select from a list of Choices with uniform probability.
type Uniform struct {
	Choices []string
}

// Choice implements the choice interface for Uniform choices.
func (u *Uniform) Choice(i uint64) (string, error) {
	choice := int(i % uint64(len(u.Choices)))
	return u.Choices[choice], nil
}

// MarshalJSON implements the json.Marshaler interface for Uniform choices.
func (u *Uniform) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Choices []string `json:"choices"`
	}{
		Choices: u.Choices,
	}
	return ffjson.Marshal(aux)
}

// Weighted is a way to select from a list of Choices with probability ratio
// supplied in Weights.
type Weighted struct {
	Choices []weightedChoice
}

type weightedChoice struct {
	name   string
	weight float64
}

func (w *weightedChoice) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Name   string  `json:"name"`
		Weight float64 `json:"weight"`
	}{
		Name:   w.name,
		Weight: w.weight,
	}
	return ffjson.Marshal(aux)
}

// Choice implements the choice interface for Weighted choices.
func (w *Weighted) Choice(i uint64) (string, error) {
	selection := make([]float64, len(w.Choices))
	cumSum := 0.0
	for ii, v := range w.Choices {
		cumSum += v.weight
		selection[ii] = cumSum
	}
	choice := uniform(i, 0, cumSum)

	for ii, v := range selection {
		if choice <= v {
			return w.Choices[ii].name, nil
		}
	}
	return "", fmt.Errorf("no selection was made")
}

// MarshalJSON implements the json.Marshaler interface for Weighted choices.
func (w *Weighted) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Choices []weightedChoice `json:"choices"`
	}{
		Choices: w.Choices,
	}
	return ffjson.Marshal(aux)
}
