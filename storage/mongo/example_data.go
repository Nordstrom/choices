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

import "github.com/foolusion/choices/storage/mongo/internal/types"

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
		&types.Namespace{
			Name:   "rands1",
			TeamID: []string{"rands"},
			Experiments: []types.Experiment{
				{
					Name:     "personalizedSort",
					Segments: allSegments,
					Params: []types.Param{
						{
							Name: "value",
							Value: types.Value{
								Choices: []string{"True", "False"},
							},
						},
					},
				},
			},
		},
		&types.Namespace{
			Name:   "rands2",
			TeamID: []string{"rands"},
			Experiments: []types.Experiment{
				{
					Name:     "categorySort",
					Segments: allSegments,
					Params: []types.Param{
						{
							Name: "value",
							Value: types.Value{
								Choices: []string{"Default", "Test1"},
							},
						},
					},
				},
			},
		},
		&types.Namespace{
			Name:   "ns1",
			TeamID: []string{"test"},
			Experiments: []types.Experiment{
				{
					Name:     "exp1",
					Segments: allSegments,
					Params: []types.Param{
						{
							Name: "buttonColor",
							Value: types.Value{
								Choices: []string{"on", "off"},
							},
						},
					},
				},
			},
		},
		&types.Namespace{
			Name:   "ns2",
			TeamID: []string{"test"},
			Experiments: []types.Experiment{
				{
					Name:     "exp2",
					Segments: allSegments,
					Params: []types.Param{
						{
							Name: "emojiSize",
							Value: types.Value{
								Choices: []string{"small", "big"},
							},
						},
						{
							Name: "emoji",
							Value: types.Value{
								Choices: []string{"💩", "😘", "😱"},
								Weights: []float64{1, 2, 3},
							},
						},
					},
				},
			},
		},
		&types.Namespace{
			Name:   "ns3",
			TeamID: []string{"test"},
			Experiments: []types.Experiment{
				{
					Name:     "exp3",
					Segments: allSegments,
					Params: []types.Param{
						{
							Name: "first",
							Value: types.Value{
								Choices: []string{"on", "off"},
							},
						},
						{
							Name: "second",
							Value: types.Value{
								Choices: []string{"on", "off"},
							},
						},
					},
				},
			},
		},
		&types.Namespace{
			Name:   "snbmow1",
			TeamID: []string{"mobilesearch"},
			Experiments: []types.Experiment{
				{
					Name:     "mobileResultTest",
					Segments: allSegments,
					Params: []types.Param{
						{
							Name: "resultCount",
							Value: types.Value{
								Choices: []string{"24", "36", "48"},
							},
						},
					},
				},
			},
		},
		&types.Namespace{
			Name:   "snb1",
			TeamID: []string{"search"},
			Experiments: []types.Experiment{
				{
					Name:     "categoryHeaderFilterTest",
					Segments: allSegments,
					Params: []types.Param{
						{
							Name: "headerExperience",
							Value: types.Value{
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
		&types.Namespace{
			Name:   "john",
			TeamID: []string{"test"},
			Experiments: []types.Experiment{
				{
					Name:     "johnHeight",
					Segments: firstHalf,
					Params: []types.Param{
						{
							Name: "height",
							Value: types.Value{
								Choices: []string{"short", "medium", "tall"},
							},
						},
					},
				},
				{
					Name:     "johnWeight",
					Segments: secondHalf,
					Params: []types.Param{
						{
							Name: "weight",
							Value: types.Value{
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
