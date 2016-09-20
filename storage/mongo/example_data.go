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

import (
	"github.com/foolusion/choices"
	"github.com/foolusion/choices/storage/mongo/internal/types"
)

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
		&types.NamespaceInput{
			Name:     "rands1",
			Segments: noSegments,
			TeamID:   []string{"rands"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "personalizedSort",
					Segments: allSegments,
					Params: []types.ParamInput{
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
		&types.NamespaceInput{
			Name:     "rands2",
			Segments: noSegments,
			TeamID:   []string{"rands"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "categorySort",
					Segments: allSegments,
					Params: []types.ParamInput{
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
		&types.NamespaceInput{
			Name:     "ns1",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "exp1",
					Segments: allSegments,
					Params: []types.ParamInput{
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
		&types.NamespaceInput{
			Name:     "ns2",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "exp2",
					Segments: allSegments,
					Params: []types.ParamInput{
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
		&types.NamespaceInput{
			Name:     "ns3",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "exp3",
					Segments: allSegments,
					Params: []types.ParamInput{
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
		&types.NamespaceInput{
			Name:     "snbmow1",
			Segments: noSegments,
			TeamID:   []string{"mobilesearch"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "mobileResultTest",
					Segments: allSegments,
					Params: []types.ParamInput{
						{
							Name: "resultCount",
							Type: choices.ValueTypeUniform,
							Value: choices.Uniform{
								Choices: []string{"24", "36", "48"},
							},
						},
					},
				},
			},
		},
		&types.NamespaceInput{
			Name:     "snb1",
			Segments: noSegments,
			TeamID:   []string{"search"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "categoryHeaderFilterTest",
					Segments: allSegments,
					Params: []types.ParamInput{
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
		&types.NamespaceInput{
			Name:     "john",
			Segments: noSegments,
			TeamID:   []string{"test"},
			Experiments: []types.ExperimentInput{
				{
					Name:     "johnHeight",
					Segments: firstHalf,
					Params: []types.ParamInput{
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
					Params: []types.ParamInput{
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
