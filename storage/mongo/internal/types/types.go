package types

// Namespace is a helper type for loading a Namespace into mongo.
type Namespace struct {
	Name        string
	TeamID      []string
	Experiments []Experiment
}

// Experiment is a helper type for loading a Experiment into mongo.
type Experiment struct {
	Name     string
	Segments string
	Params   []Param
}

// Param is a helper type for loading a Param into mongo.
type Param struct {
	Name  string
	Value Value
}

// Value is a helper type for loading a Value into mongo.
type Value struct {
	Choices []string
	Weights []float64
}
