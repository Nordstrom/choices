package types

import (
	"github.com/foolusion/choices"
	"gopkg.in/mgo.v2/bson"
)

// Namespace is a helper type to read Namespace data.
type Namespace struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string
	Segments    string
	TeamID      []string
	Experiments []Experiment
}

// Experiment is a helper type to read Experiment data.
type Experiment struct {
	Name     string
	Segments string
	Params   []Param
}

// Param is a helper type to read Param data.
type Param struct {
	Name  string
	Type  choices.ValueType
	Value bson.Raw
}

// NamespaceInput is a helper type for loading a Namespace into mongo.
type NamespaceInput struct {
	Name        string
	Segments    string
	TeamID      []string
	Experiments []ExperimentInput
}

// ExperimentInput is a helper type for loading a Experiment into mongo.
type ExperimentInput struct {
	Name     string
	Segments string
	Params   []ParamInput
}

// ParamInput is a helper type for loading a Param into mongo.
type ParamInput struct {
	Name  string
	Type  choices.ValueType
	Value interface{}
}
