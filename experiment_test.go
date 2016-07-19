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

import "testing"

func TestExperiment(t *testing.T) {
	seg := segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	tests := []struct {
		exp  Experiment
		want []paramValue
		err  error
	}{
		{
			exp: Experiment{
				Name: "experiment",
				Params: []Param{
					{Name: "p1", Value: &Uniform{Choices: []string{"a", "b"}}},
					{Name: "p2", Value: &Weighted{Choices: []string{"a", "b", "c"}, Weights: []float64{1, 10, 1}}},
				},
				Segments: seg,
			},
			want: []paramValue{{name: "p1", value: "a"}, {name: "p2", value: "c"}},
			err:  nil,
		},
	}
	h := &hashConfig{salt: "test"}
	for _, test := range tests {
		got, err := test.exp.eval(h)
		if err != test.err {
			t.Errorf("%v.eval() = %v %v, want %v %v", test.exp, got, err, test.want, test.err)
			t.FailNow()
		}
		if len(got) != len(test.want) {
			t.Errorf("%v.eval() = %v %v, want %v %v", test.exp, got, err, test.want, test.err)
			t.FailNow()
		}
		for i, v := range got {
			if v != test.want[i] {
				t.Errorf("%v.eval() = %v %v, want %v %v", test.exp, got, err, test.want, test.err)
				t.FailNow()
			}
		}
	}
}

func TestParamEval(t *testing.T) {
	tests := []struct {
		p    Param
		want paramValue
		err  error
	}{
		{
			p:    Param{Name: "test", Value: &Uniform{Choices: []string{"a", "b"}}},
			want: paramValue{name: "test", value: "a"},
			err:  nil,
		},
		{
			p:    Param{Name: "test", Value: &Weighted{Choices: []string{"a", "b"}, Weights: []float64{10, 90}}},
			want: paramValue{name: "test", value: "b"},
			err:  nil,
		},
	}
	h := &hashConfig{salt: "test"}
	for _, test := range tests {
		got, err := test.p.eval(h)
		if err != test.err {
			t.Errorf("%v.eval(nil) = %v %v, want %v %v", test.p, got, err, test.want, test.err)
			t.FailNow()
		}
		if got != test.want {
			t.Errorf("%v.eval(nil) = %v %v, want %v %v", test.p, got, err, test.want, test.err)
			t.FailNow()
		}
	}
}
