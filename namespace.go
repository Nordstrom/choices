package choices

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var config = struct {
	globalSalt string
}{
	globalSalt: "choices",
}

var expManager = struct {
	namespaceIndexByTeamID map[string][]int
	namespace              []Namespace
	mu                     sync.RWMutex
}{
	namespaceIndexByTeamID: make(map[string][]int, 100),
	namespace:              []Namespace{},
}

func nsByID(teamID string) []int {
	expManager.mu.RLock()
	defer expManager.mu.RUnlock()

	return expManager.namespaceIndexByTeamID[teamID]
}

func ns(index int) Namespace {
	expManager.mu.RLock()
	defer expManager.mu.RUnlock()

	return expManager.namespace[index]
}

type Namespace struct {
	Name        string
	Segments    [16]byte
	TeamID      []string
	Experiments []Experiment
	Units       []string
}

func (n *Namespace) eval(units map[string][]string) int {
	// TODO: need to implement the eval
	return 0
}

func filterUnits(units map[string][]string, keep []string) map[string][]string {
	out := make(map[string][]string)

	for _, v := range keep {
		out[v] = units[v]
	}
	return out
}

type Response struct{}

func (r *Response) add(i int) {}

func Namespaces(ctx context.Context, teamID string, units map[string][]string) *Response {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	response := &Response{}

	namespaces := nsByID(teamID)
	for _, index := range namespaces {
		ns := ns(index)

		u := filterUnits(units, ns.Units)

		response.add(ns.eval(u))
	}
	return response
}

type Experiment struct {
	Name       string
	Definition Definition
	Segments   [16]byte
}

type Definition struct {
	Params []Param
}

type Param struct {
	Name  string
	Type  string
	Value Value
}

type Value interface {
	Value(context.Context) string
	String() string
}

type Uniform struct {
	Choices []string
	choice  int
}

func (u *Uniform) Value(ctx context.Context) string {
	u.eval(ctx)
	return u.Choices[u.choice]
}

func (u *Uniform) String() string {
	return u.Choices[u.choice]
}

func (u *Uniform) eval(ctx context.Context) error {
	i, err := hash(ctx)
	if err != nil {
		return err
	}
	u.choice = int(i) % len(u.Choices)
	return nil
}

type Weighted struct {
	Choices []string
	Weights []float64
	choice  int
}

func (w *Weighted) Value(ctx context.Context) string {
	w.eval(ctx)
	return w.Choices[w.choice]
}

func (w *Weighted) String() string {
	return w.Choices[w.choice]
}

func (w *Weighted) eval(ctx context.Context) error {
	if len(w.Choices) != len(w.Weights) {
		return fmt.Errorf("len(w.Choices) != len(w.Weights): %v != %v", len(w.Choices), len(w.Weights))
	}

	i, err := hash(ctx)
	if err != nil {
		return err
	}

	selection := make([]float64, len(w.Weights))
	cumSum := 0.0
	for i, v := range w.Weights {
		cumSum += v
		selection[i] = cumSum
	}
	choice := getUniform(i, 0, cumSum)
	for i, v := range selection {
		if choice < v {
			w.choice = i
			return nil
		}
	}

	return fmt.Errorf("no choice made")
}
