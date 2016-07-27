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

func BenchmarkNamespaceEval(b *testing.B) {
	ns := NewNamespace("t1", "test")
	ns.AddExperiment(
		"aTest",
		[]Param{{Name: "a", Value: &Uniform{Choices: []string{"b", "c"}}}},
		128,
	)
	h := hashConfig{}
	h.setSalt(defaultSalt)
	h.setUserID("my-user-id")
	for i := 0; i < b.N; i++ {
		ns.eval(h)
	}
}
