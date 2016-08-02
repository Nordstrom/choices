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
	Segments    [16]byte
	TeamID      []string
	Experiments []MongoExperimentInput
}

type MongoExperimentInput struct {
	Name     string
	Segments [16]byte
	Params   []MongoParamInput
}

type MongoParamInput struct {
	Name  string
	Type  choices.ValueType
	Value interface{}
}

func (m *Mongo) LoadExampleData() {
	coll := m.sess.DB(m.db).C(m.coll)
	coll.RemoveAll(nil)
	coll.Insert(
		&MongoNamespaceInput{
			Name:     "ns1",
			Segments: [16]byte{},
			TeamID:   []string{"test"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "exp1",
					Segments: choices.SegmentsAll,
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
			Segments: [16]byte{},
			TeamID:   []string{"test"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "exp2",
					Segments: choices.SegmentsAll,
					Params: []MongoParamInput{
						{
							Name: "emoji",
							Type: choices.ValueTypeWeighted,
							Value: choices.Weighted{
								Choices: []string{"💩", "⬛️", "♻️"},
								Weights: []float64{1, 8, 7},
							},
						},
					},
				},
			},
		},
		&MongoNamespaceInput{
			Name:     "ns3",
			Segments: [16]byte{},
			TeamID:   []string{"test"},
			Experiments: []MongoExperimentInput{
				{
					Name:     "exp3",
					Segments: choices.SegmentsAll,
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
	)
}
