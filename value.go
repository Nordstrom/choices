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
	Value(ns, exp, p string, units []unit) string
	String() string
}

type Uniform struct {
	Choices []string
	choice  int
}

func (u *Uniform) Value(ns, exp, p string, units []unit) string {
	u.eval(ns, exp, p, units)
	return u.Choices[u.choice]
}

func (u *Uniform) String() string {
	return u.Choices[u.choice]
}

func (u *Uniform) eval(ns, exp, p string, units []unit) error {
	i, err := hash(hashNs(ns), hashExp(exp), hashParam(p), hashUnits(units))
	if err != nil {
		return err
	}
	u.choice = int(i) % len(u.Choices)
	return nil
}

type Weighted struct {
	Choices []string
	Weights []float64
	choice  int
}

func (w *Weighted) Value(ns, exp, p string, units []unit) string {
	w.eval(ns, exp, p, units)
	return w.Choices[w.choice]
}

func (w *Weighted) String() string {
	return w.Choices[w.choice]
}

func (w *Weighted) eval(ns, exp, p string, units []unit) error {
	if len(w.Choices) != len(w.Weights) {
		return fmt.Errorf(
			"len(w.Choices) != len(w.Weights): %v != %v",
			len(w.Choices),
			len(w.Weights))
	}

	i, err := hash(hashNs(ns), hashExp(exp), hashParam(p), hashUnits(units))
	if err != nil {
		return err
	}

	selection := make([]float64, len(w.Weights))
	cumSum := 0.0
	for i, v := range w.Weights {
		cumSum += v
		selection[i] = cumSum
	}
	choice := uniform(i, 0, cumSum)
	for i, v := range selection {
		if choice < v {
			w.choice = i
			return nil
		}
	}

	return fmt.Errorf("no choice made")
}
