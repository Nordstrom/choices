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

import "fmt"

// ValueType are the different types of Values a Param can have. This is used
// for parsing params in storage.
type ValueType int

// Constants for ValueTypes
const (
	ValueTypeBad ValueType = iota
	ValueTypeUniform
	ValueTypeWeighted
)

// Value is the interface Param Values must implement. They take a hash value
// and return the string that represents the value or an error.
type Value interface {
	Value(i uint64) (string, error)
}

// Uniform is a way to select from a list of Choices with uniform probability.
type Uniform struct {
	Choices []string
}

// Value implements the Value interface for Uniform choices.
func (u *Uniform) Value(i uint64) (string, error) {
	choice := int(i % uint64(len(u.Choices)))
	return u.Choices[choice], nil
}

// Weighted is a way to select from a list of Choices with probability ratio
// supplied in Weights.
type Weighted struct {
	Choices []string
	Weights []float64
}

// Value implements the Value interface for Weighted choices.
func (w *Weighted) Value(i uint64) (string, error) {
	if len(w.Choices) != len(w.Weights) {
		return "", fmt.Errorf(
			"len(w.Choices) != len(w.Weights): %v != %v",
			len(w.Choices),
			len(w.Weights))
	}

	selection := make([]float64, len(w.Weights))
	cumSum := 0.0
	for ii, v := range w.Weights {
		cumSum += v
		selection[ii] = cumSum
	}
	choice := uniform(i, 0, cumSum)

	for ii, v := range selection {
		if choice <= v {
			return w.Choices[ii], nil
		}
	}
	return "", fmt.Errorf("no selection was made")
}
