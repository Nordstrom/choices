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
)

func TestFromExperiment(t *testing.T) {
	tests := map[string]struct {
		in   *storage.Experiment
		want Experiment
		err  error
	}{
		"emptyExperiment": {in: &storage.Experiment{}, want: Experiment{}, err: nil},
		"oneExperiment": {
			in: &storage.Experiment{
				Namespace: "ns",
				Name:      "exp1",
				Labels:    map[string]string{"team": "test"},
				Segments:  []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
				Params:    []*storage.Param{{Name: "param1", Value: &storage.Value{Choices: []string{"a", "b"}}}},
			},
			want: Experiment{
				Namespace: "ns",
				Name:      "exp1",
				Labels:    map[string]string{"team": "test"},
				Segments:  segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
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
		if out.Segments != test.want.Segments {
			t.Errorf("%s segments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if len(out.Params) != len(test.want.Params) {
			t.Errorf("%s experiments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
	}
}
