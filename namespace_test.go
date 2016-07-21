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
	"log"
	"testing"
)

func TestNsByID(t *testing.T) {
}

func BenchmarkNamespaces(b *testing.B) {
	ns, err := NewNamespace("t1", "test", []string{"userid"})
	if err != nil {
		log.Fatal("%v", err)
	}
	ns.Addexp(
		"aTest",
		[]Param{{Name: "a", Value: &Uniform{Choices: []string{"b", "c"}}}},
		128,
	)
	if err := Addns(ns); err != nil {
		log.Fatalf("%v", err)
	}

	teamID := "test"
	units := map[string][]string{"userid": []string{"some-user-id"}}
	for i := 0; i < b.N; i++ {
		Namespaces(nil, teamID, units)
	}
}

func BenchmarkNamespaceEval(b *testing.B) {
	ns, err := NewNamespace("t1", "test", []string{"userid"})
	if err != nil {
		log.Fatal("%v", err)
	}
	ns.Addexp(
		"aTest",
		[]Param{{Name: "a", Value: &Uniform{Choices: []string{"b", "c"}}}},
		128,
	)
	units := []unit{{key: "userid", value: []string{"my-super-unique-userid"}}}
	for i := 0; i < b.N; i++ {
		ns.eval(units)
	}
}

func BenchmarkFilterUnits(b *testing.B) {
	units := map[string][]string{"test": []string{"a", "b"}, "useless": []string{"blah"}, "keep": []string{"arst"}}
	keep := []string{"keep", "test"}
	for i := 0; i < b.N; i++ {
		filterUnits(units, keep)
	}
}
