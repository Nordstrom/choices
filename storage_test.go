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

import (
	"testing"

	"github.com/foolusion/elwinprotos/storage"
)

func TestFromNamespace(t *testing.T) {
	tests := map[string]struct {
		in   *storage.Namespace
		want Namespace
		err  error
	}{
		"emptyNamespace": {in: &storage.Namespace{}, want: Namespace{}, err: nil},
		"oneExperiment":  {in: &storage.Namespace{Name: "ns", Labels: []string{"test"}, Experiments: []*storage.Experiment{{Name: "exp1", Segments: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, Params: []*storage.Param{{Name: "param1", Value: &storage.Value{Choices: []string{"a", "b"}}}}}}}, want: Namespace{Name: "ns", Labels: []string{"test"}, Segments: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, Experiments: []Experiment{{Name: "exp1", Segments: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, Params: []Param{{Name: "param1", Value: &Uniform{Choices: []string{"a", "b"}}}}}}}, err: nil},
	}

	for k, test := range tests {
		out, err := FromNamespace(test.in)
		if out.Name != test.want.Name {
			t.Errorf("%s name: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if out.Segments != test.want.Segments {
			t.Errorf("%s segments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if len(out.Labels) != len(test.want.Labels) {
			t.Errorf("%s teamID: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if len(out.Experiments) != len(test.want.Experiments) {
			t.Errorf("%s experiments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if err != test.err {
			t.Errorf("%s: FromNamespace(%+v) = %+v want %+v", k, test.in, err, test.err)
		}
	}
}
