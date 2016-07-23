package choices

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongo struct{}

type MongoNamespace struct {
	Name        string
	Segments    [16]byte
	TeamID      []string
	Experiments []MongoExperiments
}

type MongoExperiments struct {
	Name     string
	Segments [16]byte
	Params   []MongoParams
}

type MongoParams struct {
	Name  string
	Type  ValueType
	Value bson.Raw
}

func (m *mongo) TeamNamespaces(teamID string) []Namespace {
	sess, err := mgo.Dial("elwin-storage")
	if err != nil {
		log.Println(err)
	}
	c := sess.DB("test").C("namespaces")
	iter := c.Find(bson.M{"teamID": teamID}).Iter()
	var mongoNamespaces []MongoNamespace
	err = iter.All(&mongoNamespaces)
	if err != nil {
		log.Println(err)
	}

	namespaces := make([]Namespace, len(mongoNamespaces))
	for i, n := range mongoNamespaces {
		namespaces[i] = Namespace{Name: n.Name, Segments: n.Segments, TeamID: n.TeamID, Experiments: make([]Experiment, len(n.Experiments))}
		for j, e := range n.Experiments {
			namespaces[i].Experiments[j] = Experiment{Name: e.Name, Segments: e.Segments, Params: make([]Param, len(e.Params))}
			for k, p := range e.Params {
				namespaces[i].Experiments[j].Params[k] = Param{Name: p.Name}
				switch p.Type {
				case ValueTypeUniform:
					var uniform Uniform
					p.Value.Unmarshal(&uniform)
					namespaces[i].Experiments[j].Params[k].Value = &uniform
				case ValueTypeWeighted:
					var weighted Weighted
					p.Value.Unmarshal(&weighted)
					namespaces[i].Experiments[j].Params[k].Value = &weighted
				}
			}
		}
	}
	return namespaces
}
