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
	"testing"

	"github.com/foolusion/elwinprotos/storage"
	"k8s.io/apimachinery/pkg/labels"
)

func TestExperiment(t *testing.T) {
	backup := globalSalt
	defer func() {
		globalSalt = backup
	}()
	globalSalt = ""
	seg := make(segments, 16)
	copy(seg[:], segmentsAll[:])
	tests := []struct {
		exp  Experiment
		want ExperimentResponse
		err  error
	}{
		{
			exp: Experiment{
				Name: "experiment",
				Params: []Param{
					{Name: "p1", Choices: &Uniform{Choices: []string{"a", "b"}}},
					{Name: "p2", Choices: &Weighted{Choices: []weightedChoice{{"a", 1}, {"b", 10}, {"c", 1}}}},
				},
				Segments: seg,
			},
			want: ExperimentResponse{Name: "experiment", Namespace: "", Params: []ParamValue{{Name: "p1", Value: "b"}, {Name: "p2", Value: "b"}}},
			err:  nil,
		},
	}
	h := hashConfig{}
	for _, test := range tests {
		got, err := test.exp.eval(h)
		if err != test.err {
			t.Fatalf("%v.eval() = %v %v, want %v %v", test.exp, got, err, test.want, test.err)
		}
		if test.want.Name != got.Name || test.want.Namespace != got.Namespace {
			t.Fatalf("%v.eval() = %v, want %v", test.exp, test.want, got)
		}
		for i, v := range got.Params {
			if v != test.want.Params[i] {
				t.Fatalf("%v.eval() = %v %v, want %v %v", test.exp, got, err, test.want, test.err)
			}
		}
	}
}

func TestParamEval(t *testing.T) {
	backup := globalSalt
	defer func() {
		globalSalt = backup
	}()

	globalSalt = "test"
	tests := []struct {
		p    Param
		want ParamValue
		err  error
	}{
		{
			p:    Param{Name: "test", Choices: &Uniform{Choices: []string{"a", "b"}}},
			want: ParamValue{Name: "test", Value: "b"},
			err:  nil,
		},
		{
			p:    Param{Name: "test", Choices: &Weighted{Choices: []weightedChoice{{"a", 10}, {"b", 90}}}},
			want: ParamValue{Name: "test", Value: "b"},
			err:  nil,
		},
	}
	h := hashConfig{salt: [3]string{"", "", ""}}
	for _, test := range tests {
		got, err := test.p.eval(h)
		if err != test.err {
			t.Errorf("%v.eval(nil) = %v %v, want %v %v", test.p, got, err, test.want, test.err)
		}
		if got != test.want {
			t.Errorf("%v.eval(nil) = %v %v, want %v %v", test.p, got, err, test.want, test.err)
		}
	}
}

func BenchmarkExperimentEval(b *testing.B) {
	e := Experiment{
		Name:      "experiment",
		Namespace: "foo",
		Segments:  segmentsAll,
		Params: []Param{
			{Name: "p", Choices: &Uniform{Choices: []string{"a", "b"}}},
		},
		Labels: map[string]string{"test": "true"},
	}
	h := hashConfig{
		salt: [3]string{"namespace", "", ""},
	}
	for i := 0; i < b.N; i++ {
		if _, err := e.eval(h); err != nil {
			b.Fatal(err)
		}
	}
}

func TestSetSegments(t *testing.T) {
	e := NewExperiment("test")
	tests := map[string]struct {
		seg  segments
		want segments
	}{
		"none": {seg: segments{}, want: segments{}},
		"all":  {seg: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, want: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}},
		"some": {seg: segments{255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0}, want: segments{255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0}},
	}
	for tname, test := range tests {
		e.SetSegments(test.seg)
		for i := range e.Segments {
			if e.Segments[i] != test.want[i] {
				t.Fatalf("%s: e.SetSegments(%v) = %v, want %v", tname, test.seg, e.Segments, test.want)
			}
		}
	}
}

func TestToExperiment(t *testing.T) {
	tests := map[string]struct {
		e    Experiment
		want storage.Experiment
	}{
		"simple": {
			e:    Experiment{Name: "test", Labels: labels.Set{"team": "ato", "platform": "desktop"}, Segments: segments{}, Params: []Param{}},
			want: storage.Experiment{Name: "test", Labels: map[string]string{"team": "ato", "platform": "desktop"}, Segments: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},
	}
	for tname, test := range tests {
		out := test.e.ToExperiment()
		if out.Name != test.want.Name {
			t.Fatalf("%s: e.ToExperiment() = %v, want %v", tname, *out, test.want)
		}
		for k, v := range test.want.Labels {
			if ov, ok := out.Labels[k]; !ok {
				t.Fatalf("%s: e.ToExperiment(): key %v does not exist in output", tname, k)
			} else if v != ov {
				t.Fatalf("%s: e.ToExperiment(): key %v = got %v, want %v", tname, k, ov, v)
			}
		}
	}
}
