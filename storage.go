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
	"time"

	"golang.org/x/net/context"
)

// Storage is an interface for the storage of experiments. Storage has two
// functions. Update, which checks for new data from the data store, and Read
// which returns the current slice of Namespaces. Clients using the storage
// interface should never write data to the slice returned by Namespaces.
// Storage should be read only the values should never be overwritten.
type Storage interface {
	Update()
	Read() []Namespace
}

// TeamNamespaces filters the namespaces from storage based on teamID.
func TeamNamespaces(s Storage, teamID string) []Namespace {
	allNamespaces := s.Read()
	teamNamespaces := make([]Namespace, 0, len(allNamespaces))
	for _, n := range allNamespaces {
		for _, t := range n.TeamID {
			if t == teamID {
				teamNamespaces = append(teamNamespaces, n)
			}
		}
	}
	return teamNamespaces
}

// SetStorage sets the storage engine. It starts a ticker that will call
// s.Update() until the context is cancelled. To change the tick interval call
// SetUpdateInterval(d time.Duration). Must cancel the context before calling
// SetStorage agian otherwise you will leak go routines.
func SetStorage(ctx context.Context, s Storage) {
	config.storage = s
	go func() {
		s.Update()
		c := time.Tick(config.updateInterval)
		for {
			select {
			case <-c:
				s.Update()
			case <-ctx.Done():
				return
			}
		}
	}()
}
