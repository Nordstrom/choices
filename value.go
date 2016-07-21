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

type Value interface {
	Value(i uint64) (string, error)
	String() string
}

type Uniform struct {
	Choices []string
	choice  int
}

func (u *Uniform) Value(i uint64) (string, error) {
	u.choice = int(i % uint64(len(u.Choices)))
	return u.Choices[u.choice], nil
}

func (u *Uniform) String() string {
	return u.Choices[u.choice]
}

type Weighted struct {
	Choices []string
	Weights []float64
	choice  int
}

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
	selected := false
	for ii, v := range selection {
		if choice < v {
			w.choice = ii
			selected = true
			break
		}
	}

	if !selected {
		return "", fmt.Errorf("no selection was made")
	}

	return w.Choices[w.choice], nil
}

func (w *Weighted) String() string {
	return w.Choices[w.choice]
}
