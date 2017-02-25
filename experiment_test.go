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

func TestExperiment(t *testing.T) {
	backup := globalSalt
	defer func() {
		globalSalt = backup
	}()
	globalSalt = ""
	var seg segments
	copy(seg[:], segmentsAll[:])
	tests := []struct {
		exp  Experiment
		want []ParamValue
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
			want: []ParamValue{{Name: "p1", Value: "b"}, {Name: "p2", Value: "b"}},
			err:  nil,
		},
	}
	h := hashConfig{}
	for _, test := range tests {
		got, err := test.exp.eval(h)
		if err != test.err {
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

func TestExperimentSampleSegments(t *testing.T) {
	tests := map[string]struct {
		nsSeg   segments
		num     int
		nsWant  segments
		expWant segments
	}{
		"all": {
			nsSeg:   segments{},
			num:     128,
			nsWant:  segmentsAll,
			expWant: segmentsAll,
		},
		"half": {
			nsSeg:   segments{255, 255, 255, 255, 255, 255, 255, 255},
			num:     64,
			nsWant:  segmentsAll,
			expWant: segments{0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		"too much": {
			nsSeg:   segments{},
			num:     9000,
			nsWant:  segmentsAll,
			expWant: segmentsAll,
		},
	}

	ns := NewNamespace("test")

	for k, test := range tests {
		ns.Segments = test.nsSeg
		e := NewExperiment("e")
		e = e.SampleSegments(ns, test.num)
		if e.Segments != test.expWant {
			t.Errorf("%s: experient segments: %v, want %v", k, e.Segments, test.expWant)
		}
		if ns.Segments != test.nsWant {
			t.Errorf("%s: namespace segments: %v, want %v", k, ns.Segments, test.nsWant)
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
			p:    Param{Name: "test", Value: &Uniform{Choices: []string{"a", "b"}}},
			want: ParamValue{Name: "test", Value: "b"},
			err:  nil,
		},
		{
			p:    Param{Name: "test", Value: &Weighted{Choices: []string{"a", "b"}, Weights: []float64{10, 90}}},
			want: ParamValue{Name: "test", Value: "b"},
			err:  nil,
		},
	}
	h := hashConfig{salt: [3]string{"", "", ""}}
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

func BenchmarkExperimentEval(b *testing.B) {
	e := Experiment{
		Name: "experiment",
		Params: []Param{
			{Name: "p", Value: &Uniform{Choices: []string{"a", "b"}}},
		},
	}
	copy(e.Segments[:], segmentsAll[:])
	h := hashConfig{
		salt: [3]string{"namespace", "", ""},
	}
	for i := 0; i < b.N; i++ {
		if _, err := e.eval(h); err != nil {
			b.Fatal(err)
		}
	}
}
