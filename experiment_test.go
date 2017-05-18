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

	"github.com/Nordstrom/choices/util"
	"github.com/foolusion/elwinprotos/storage"
	"k8s.io/apimachinery/pkg/labels"
)

func TestExperiment(t *testing.T) {
	backup := globalSalt
	defer func() {
		globalSalt = backup
	}()
	globalSalt = ""
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
				Segments: &segmentsAll,
			},
			want: ExperimentResponse{Name: "experiment", Namespace: "", Params: []*ParamValue{{Name: "p1", Value: "b"}, {Name: "p2", Value: "b"}}},
			err:  nil,
		},
	}
	h := hashConfig{}
	for _, test := range tests {
		got := new(ExperimentResponse)
		err := test.exp.eval(got, h)
		if err != test.err {
			t.Fatalf("%v.eval() = %v %v, want %v %v", test.exp, got, err, test.want, test.err)
		}
		if test.want.Name != got.Name || test.want.Namespace != got.Namespace {
			t.Fatalf("%v.eval() = %v, want %v", test.exp, test.want, got)
		}
		for i, v := range got.Params {
			if *v != *test.want.Params[i] {
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
		got := new(ParamValue)
		err := test.p.eval(got, h)
		if err != test.err {
			t.Errorf("%v.eval(nil) = %v %v, want %v %v", test.p, got, err, test.want, test.err)
		}
		if *got != test.want {
			t.Errorf("%v.eval(nil) = %v %v, want %v %v", test.p, got, err, test.want, test.err)
		}
	}
}

func BenchmarkExperimentEval(b *testing.B) {
	e := Experiment{
		Name:      "experiment",
		Namespace: "foo",
		Segments:  &segmentsAll,
		Params: []Param{
			{Name: "p", Choices: &Uniform{Choices: []string{"a", "b"}}},
		},
		Labels: map[string]string{"test": "true"},
	}
	h := hashConfig{
		salt: [3]string{"namespace", "", ""},
	}
	er := new(ExperimentResponse)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := e.eval(er, h); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParamEval(b *testing.B) {
	p := Param{
		Name: "param",
		Choices: &Uniform{
			Choices: []string{
				"a",
				"b",
			},
		},
	}
	h := hashConfig{
		salt:   [3]string{"namespace", "experiment", ""},
		userID: "andrew",
	}
	ep := new(ParamValue)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := p.eval(ep, h); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark1Experiments(b *testing.B) {
	exps := make([]*Experiment, 1)
	for i := range exps {
		exps[i] = &Experiment{
			Name:      util.BasicNameGenerator.GenerateName("e-"),
			Namespace: util.BasicNameGenerator.GenerateName("ns-"),
			Labels: labels.Set{
				"team":     "ato",
				"platform": "service",
			},
			Segments: &segments{
				b:   []byte{255, 255, 255, 255, 255, 255, 255, 255},
				len: 8,
			},
			Params: []Param{
				{
					Name: util.BasicNameGenerator.GenerateName("p-"),
					Choices: &Uniform{
						Choices: []string{
							util.BasicNameGenerator.GenerateName("c-"),
							util.BasicNameGenerator.GenerateName("c-"),
						},
					},
				},
			},
		}
	}
	expStore := &experimentStore{
		cache: exps,
	}

	sel, err := labels.Parse("team in (ato)")
	if err != nil {
		b.Error(err)
	}

	c := Config{
		storage: expStore,
	}
	resp := make([]*ExperimentResponse, len(exps))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err = c.Experiments(resp, "andrew", sel)
		if err != nil {
			b.Error(err)
		}
	}
}

func randTeam(i int) string {
	teams := []string{"ato", "epe"}
	return teams[i%len(teams)]
}

func Benchmark100Experiments(b *testing.B) {
	exps := make([]*Experiment, 100)
	for i := range exps {
		exps[i] = &Experiment{
			Name:      util.BasicNameGenerator.GenerateName("e-"),
			Namespace: util.BasicNameGenerator.GenerateName("ns-"),
			Labels: labels.Set{
				"team":     "ato",
				"platform": "service",
			},
			Segments: &segments{
				b:   []byte{255, 255, 255, 255, 255, 255, 255, 255},
				len: 8,
			},
			Params: []Param{
				{
					Name: util.BasicNameGenerator.GenerateName("p-"),
					Choices: &Uniform{
						Choices: []string{
							util.BasicNameGenerator.GenerateName("c-"),
							util.BasicNameGenerator.GenerateName("c-"),
						},
					},
				},
			},
		}
	}
	expStore := &experimentStore{
		cache: exps,
	}

	sel, err := labels.Parse("team in (ato)")
	if err != nil {
		b.Error(err)
	}

	c := Config{
		storage: expStore,
	}
	resp := make([]*ExperimentResponse, len(exps))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err = c.Experiments(resp, "andrew", sel)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestToExperiment(t *testing.T) {
	tests := map[string]struct {
		e    Experiment
		want storage.Experiment
	}{
		"simple": {
			e:    Experiment{Name: "test", Labels: labels.Set{"team": "ato", "platform": "desktop"}, Segments: &segments{}, Params: []Param{}},
			want: storage.Experiment{Name: "test", Labels: map[string]string{"team": "ato", "platform": "desktop"}, Segments: &storage.Segments{}},
		},
	}
	for tname, test := range tests {
		out := test.e.ToExperiment(nil)
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
