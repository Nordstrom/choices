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

package mongo

import "github.com/foolusion/choices"

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

const (
	noSegments  = "00000000000000000000000000000000"
	firstHalf   = "ffffffffffffffff0000000000000000"
	secondHalf  = "0000000000000000ffffffffffffffff"
	allSegments = "ffffffffffffffffffffffffffffffff"
)

// LoadExampleData loads the test data into the database.
func (m *Mongo) LoadExampleData() {
	coll := m.sess.DB(m.db).C(m.coll)
	coll.RemoveAll(nil)
	coll.Insert(
		&NamespaceInput{
			Name:     "rands1",
			Segments: noSegments,
			TeamID:   []string{"rands"},
			Experiments: []ExperimentInput{
				{
					Name:     "personalizedSort",
					Segments: allSegments,
					Params: []ParamInput{
						{
							Name: "value",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"True", "False"},
							},
						},
					},
				},
			},
		},
		&NamespaceInput{
			Name:     "rands2",
			Segments: noSegments,
			TeamID:   []string{"rands"},
			Experiments: []ExperimentInput{
				{
					Name:     "categorySort",
					Segments: allSegments,
					Params: []ParamInput{
						{
							Name: "value",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"Default", "Test1"},
							},
						},
					},
				},
			},
		},
		&NamespaceInput{
			Name:     "ns1",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []ExperimentInput{
				{
					Name:     "exp1",
					Segments: allSegments,
					Params: []ParamInput{
						{
							Name: "buttonColor",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"on", "off"},
							},
						},
					},
				},
			},
		},
		&NamespaceInput{
			Name:     "ns2",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []ExperimentInput{
				{
					Name:     "exp2",
					Segments: allSegments,
					Params: []ParamInput{
						{
							Name: "emojiSize",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"small", "big"},
							},
						},
						{
							Name: "emoji",
							Type: choices.ValueTypeWeighted,
							Value: choices.Weighted{
								Choices: []string{"💩", "😘", "😱"},
								Weights: []float64{1, 2, 3},
							},
						},
					},
				},
			},
		},
		&NamespaceInput{
			Name:     "ns3",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []ExperimentInput{
				{
					Name:     "exp3",
					Segments: allSegments,
					Params: []ParamInput{
						{
							Name: "first",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"on", "off"},
							},
						},
						{
							Name: "second",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"on", "off"},
							},
						},
					},
				},
			},
		},
		&NamespaceInput{
			Name:     "snbmow1",
			Segments: noSegments,
			TeamID:   []string{"mobilesearch"},
			Experiments: []ExperimentInput{
				{
					Name:     "mobileResultTest",
					Segments: allSegments,
					Params: []ParamInput{
						{
							Name: "resultCount",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"24", "48", "72"},
							},
						},
					},
				},
			},
		},
		&NamespaceInput{
			Name:     "snb1",
			Segments: noSegments,
			TeamID:   []string{"search"},
			Experiments: []ExperimentInput{
				{
					Name:     "categoryHeaderFilterTest",
					Segments: allSegments,
					Params: []ParamInput{
						{
							Name: "headerExperience",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{
									"control",
									"suppressHeaders",
									"alignNav",
								},
							},
						},
					},
				},
			},
		},
		&NamespaceInput{
			Name:     "john",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []ExperimentInput{
				{
					Name:     "johnHeight",
					Segments: firstHalf,
					Params: []ParamInput{
						{
							Name: "height",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"short", "medium", "tall"},
							},
						},
					},
				},
				{
					Name:     "johnWeight",
					Segments: secondHalf,
					Params: []ParamInput{
						{
							Name: "weight",
							Type: choices.ValueTypeUniform,
							Value: choices.Weighted{
								Choices: []string{"skinny", "average", "300"},
								Weights: []float64{1, 2, 3},
							},
						},
					},
				},
			},
		},
	)
}
