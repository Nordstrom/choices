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
	"context"
	"fmt"
	"time"
)

type ElwinConfig struct {
	globalSalt     string
	Storage        Storage
	updateInterval time.Duration
}

var config = ElwinConfig{
	globalSalt:     "choices",
	updateInterval: 5 * time.Minute,
}

// NewElwin sets the storage engine. It starts a ticker that will call
// s.Update() until the context is cancelled. To change the tick interval call
// SetUpdateInterval(d time.Duration). Must cancel the context before calling
// NewElwin again otherwise you will leak go routines.
func NewElwin(ctx context.Context, opts ...func(*ElwinConfig) error) (*ElwinConfig, error) {
	e := &ElwinConfig{
		globalSalt:     "choices",
		updateInterval: 5 * time.Minute,
	}
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	if e.Storage == nil {
		return nil, fmt.Errorf("must supply a storage option")
	}

	go func() {
		e.Storage.Update()
		c := time.Tick(config.updateInterval)
		for {
			select {
			case <-c:
				e.Storage.Update()
			case <-ctx.Done():
				return
			}
		}
	}()
	return e, nil
}

// WithGlobalSalt sets the salt used in hashing users.
func WithGlobalSalt(salt string) func(*ElwinConfig) error {
	return func(ec *ElwinConfig) error {
		ec.globalSalt = salt
		return nil
	}
}

// UpdateInterval changes the update interval for Storage. Must call
// SetStorage after this or cancel context of the current Storage and call
// SetStorage again.
func UpdateInterval(dur time.Duration) func(*ElwinConfig) error {
	return func(ec *ElwinConfig) error {
		ec.updateInterval = dur
		return nil
	}
}
