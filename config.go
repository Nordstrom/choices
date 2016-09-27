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

	"github.com/pkg/errors"

	"golang.org/x/net/context"
)

// Config is the configuration struct used in an elwin server.
type Config struct {
	globalSalt     string
	Storage        *NamespaceStore
	updateInterval time.Duration
	ErrChan        chan error
}

type ConfigOpt func(*Config) error

const (
	defaultSalt           string        = "choices"
	defaultUpdateInterval time.Duration = 5 * time.Minute
)

// NewChoices sets the storage engine. It starts a ticker that will call
// s.Update() until the context is cancelled. To change the tick interval call
// SetUpdateInterval(d time.Duration). Must cancel the context before calling
// NewChoices again otherwise you will leak go routines.
func NewChoices(ctx context.Context, opts ...ConfigOpt) (*Config, error) {
	e := &Config{
		globalSalt:     defaultSalt,
		updateInterval: defaultUpdateInterval,
		ErrChan:        make(chan error, 1),
	}
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	go func(e *Config) {
		err := e.Storage.Update()
		if err != nil {
			e.ErrChan <- errors.Wrap(err, "could not update strorage")
		}
		ticker := time.NewTicker(e.updateInterval)
		for {
			select {
			case <-ticker.C:
				err := e.Storage.Update()
				if err != nil {
					e.ErrChan <- err
					ticker.Stop()
					return
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}(e)
	return e, nil
}

// GlobalSalt sets the salt used in hashing users.
func GlobalSalt(salt string) ConfigOpt {
	return func(ec *Config) error {
		ec.globalSalt = salt
		return nil
	}
}

// UpdateInterval changes the update interval for Storage. Must call
// SetStorage after this or cancel context of the current Storage and call
// SetStorage again.
func UpdateInterval(dur time.Duration) ConfigOpt {
	return func(ec *Config) error {
		ec.updateInterval = dur
		return nil
	}
}
