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
							Name: "excrementColor",
							Type: choices.ValueTypeWeighted,
							Value: choices.Weighted{
								Choices: []string{"brown", "black", "green"},
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
