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

type MongoNamespaceInput struct {
	Name        string
	Segments    string
	TeamID      []string
	Experiments []MongoExperimentInput
}

type MongoExperimentInput struct {
	Name     string
	Segments string
	Params   []MongoParamInput
}

type MongoParamInput struct {
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

func (m *Mongo) LoadExampleData() {
	coll := m.sess.DB(m.db).C(m.coll)
	coll.RemoveAll(nil)
	coll.Insert(
		&MongoNamespaceInput{
			Name:     "rands1",
			Segments: noSegments,
			TeamID:   []string{"rands"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "personalizedSort",
					Segments: allSegments,
					Params: []MongoParamInput{
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
		&MongoNamespaceInput{
			Name:     "rands2",
			Segments: noSegments,
			TeamID:   []string{"rands"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "categorySort",
					Segments: allSegments,
					Params: []MongoParamInput{
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
		&MongoNamespaceInput{
			Name:     "ns1",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "exp1",
					Segments: allSegments,
					Params: []MongoParamInput{
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
		&MongoNamespaceInput{
			Name:     "ns2",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "exp2",
					Segments: allSegments,
					Params: []MongoParamInput{
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
		&MongoNamespaceInput{
			Name:     "ns3",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "exp3",
					Segments: allSegments,
					Params: []MongoParamInput{
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
		&MongoNamespaceInput{
			Name:     "snb2",
			Segments: noSegments,
			TeamID:   []string{"search"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "searchResultTest",
					Segments: allSegments,
					Params: []MongoParamInput{
						{
							Name: "resultCount",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"66", "99"},
							},
						},
					},
				},
			},
		},
		&MongoNamespaceInput{
			Name:     "snbmow1",
			Segments: noSegments,
			TeamID:   []string{"mobilesearch"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "mobileResultTest",
					Segments: allSegments,
					Params: []MongoParamInput{
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
		&MongoNamespaceInput{
			Name:     "snb1",
			Segments: noSegments,
			TeamID:   []string{"search"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "categoryHeaderFilterTest",
					Segments: allSegments,
					Params: []MongoParamInput{
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
		&MongoNamespaceInput{
			Name:     "john",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "johnHeight",
					Segments: firstHalf,
					Params: []MongoParamInput{
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
					Params: []MongoParamInput{
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
