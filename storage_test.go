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

	"k8s.io/apimachinery/pkg/labels"

	"github.com/Nordstrom/choices/util"
	"github.com/foolusion/elwinprotos/storage"
)

func TestFromExperiment(t *testing.T) {
	tests := map[string]struct {
		in   *storage.Experiment
		want Experiment
		err  error
	}{
		"emptyExperiment": {in: &storage.Experiment{}, want: Experiment{Segments: &segments{make([]byte, 16), 128}}, err: nil},
		"oneExperiment": {
			in: &storage.Experiment{
				Namespace: "ns",
				Name:      "exp1",
				Labels:    map[string]string{"team": "test"},
				Segments:  &storage.Segments{Len: 128, B: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}},
				Params:    []*storage.Param{{Name: "param1", Value: &storage.Value{Choices: []string{"a", "b"}}}},
			},
			want: Experiment{
				Namespace: "ns",
				Name:      "exp1",
				Labels:    map[string]string{"team": "test"},
				Segments:  &segments{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, 128},
				Params:    []Param{{Name: "param1", Choices: &Uniform{Choices: []string{"a", "b"}}}},
			},
			err: nil,
		},
	}

	for k, test := range tests {
		out := FromExperiment(test.in)
		if out.Name != test.want.Name {
			t.Errorf("%s name: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		for i := range out.Segments.b {
			if out.Segments.b[i] != test.want.Segments.b[i] {
				t.Errorf("%s segments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
			}
		}
		if len(out.Params) != len(test.want.Params) {
			t.Errorf("%s experiments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
	}
}

func BenchmarkTeamNamespaces(b *testing.B) {
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

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		teamNamespaces(expStore, sel)
	}
}
