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

package choicestest

import (
	"testing"

	"github.com/foolusion/choices"
	"github.com/foolusion/choices/storage/mem"

	"golang.org/x/net/context"
)

func BenchmarkNamespaces(b *testing.B) {
	ms := &mem.MemStore{}
	ns := choices.NewNamespace("t1", "test")
	ns.AddExperiment(
		"aTest",
		[]choices.Param{{Name: "a", Value: &choices.Uniform{Choices: []string{"b", "c"}}}},
		128,
	)
	ms.AddNamespace(ns)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	elwin, err := choices.NewChoices(ctx, mem.WithMemStore(ms))
	if err != nil {
		b.Fatalf("elwin err: %v", err)
		return
	}
	teamID := "test"
	userID := "my-user-id"
	for i := 0; i < b.N; i++ {
		elwin.Namespaces(teamID, userID)
	}
}
