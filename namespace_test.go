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

import "testing"

func TestNamespaceEval(t *testing.T) {
	tests := map[string]struct {
		ns   Namespace
		want ExperimentResponse
		err  error
	}{
		"simple": {
			ns:   Namespace{Name: "test", Segments: segmentsAll, Experiments: []Experiment{{Name: "simple", Segments: segmentsAll, Params: []Param{{Name: "p1", Choices: &Uniform{Choices: []string{"a", "b"}}}}}}},
			want: ExperimentResponse{Name: "simple", Namespace: "test", Params: []ParamValue{{Name: "p1", Value: "a"}}},
			err:  nil,
		},
		"none": {
			ns:   Namespace{Name: "2", Segments: segments{}},
			want: ExperimentResponse{},
			err:  ErrSegmentNotInExperiment,
		},
	}
	h := hashConfig{}
	h.setUserID("11171986")
	for tname, test := range tests {
		if er, err := test.ns.eval(h); err != test.err {
			t.Fatalf("%s: err = %v, want %v", tname, err, test.err)
		} else if er.Name != test.want.Name {
			t.Fatalf("%s: name = %v, want %v", tname, er.Name, test.want.Name)
		} else if er.Namespace != test.want.Namespace {
			t.Fatalf("%s: namespace = %v, want %v", tname, er.Namespace, test.want.Namespace)
		}
	}
}

func BenchmarkNamespaceEval(b *testing.B) {
	ns := NewNamespace("t1")
	e := NewExperiment("aTest").SetSegments(segmentsAll)
	e.Params = []Param{{Name: "a", Choices: &Uniform{Choices: []string{"b", "c"}}}}
	if err := ns.addExperiment(*e); err != nil {
		b.Fatal(err)
	}
	h := hashConfig{}
	h.setUserID("my-user-id")
	for i := 0; i < b.N; i++ {
		if _, err := ns.eval(h); err != nil {
			b.Fatal(err)
		}
	}
}
