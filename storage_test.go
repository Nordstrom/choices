package choices

import (
	"testing"

	storage "github.com/foolusion/choices/elwinstorage"
)

func TestFromNamespace(t *testing.T) {
	tests := map[string]struct {
		in   *storage.Namespace
		want Namespace
		err  error
	}{
		"emptyNamespace": {in: &storage.Namespace{}, want: Namespace{Segments: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}}, err: nil},
		"oneExperiment":  {in: &storage.Namespace{Name: "ns", Labels: []string{"test"}, Experiments: []*storage.Experiment{{Name: "exp1", Segments: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, Params: []*storage.Param{{Name: "param1", Value: &storage.Value{Choices: []string{"a", "b"}}}}}}}, want: Namespace{Name: "ns", TeamID: []string{"test"}, Experiments: []Experiment{{Name: "exp1", Segments: segments{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, Params: []Param{{Name: "param1", Value: &Uniform{Choices: []string{"a", "b"}}}}}}}, err: nil},
	}

	for k, test := range tests {
		out, err := FromNamespace(test.in)
		if out.Name != test.want.Name {
			t.Errorf("%s name: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if out.Segments != test.want.Segments {
			t.Errorf("%s segments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if len(out.TeamID) != len(test.want.TeamID) {
			t.Errorf("%s teamID: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if len(out.Experiments) != len(test.want.Experiments) {
			t.Errorf("%s experiments: FromNamespace(%+v) = %+v want %+v", k, test.in, out, test.want)
		}
		if err != test.err {
			t.Errorf("%s: FromNamespace(%+v) = %+v want %+v", k, err, test.err)
		}
	}
}
