package choices

import (
	"github.com/Nordstrom/choices/util"
	"github.com/foolusion/elwinprotos/storage"
)

type Namespace struct {
	Name        string
	NumSegments int
	Segments    *segments
}

func (n *Namespace) ToNamespace(ns *storage.Namespace) *storage.Namespace {
	if ns == nil {
		ns = new(storage.Namespace)
	}
	ns.Name = n.Name
	ns.NumSegments = int64(n.NumSegments)
	if ns.Segments == nil {
		ns.Segments = new(storage.Segments)
	}
	ns.Segments.B = n.Segments.b
	ns.Segments.Len = int64(n.Segments.len)
	return ns
}

const defaultNumSegments = 128

func newNamespace(name string, numSegments int) *Namespace {
	if name == "" {
		name = util.BasicNameGenerator.GenerateName("ns-")
	}
	if numSegments <= 0 {
		numSegments = defaultNumSegments
	}
	numBytes := numSegments / 8
	if numSegments%8 != 0 {
		numBytes++
	}
	return &Namespace{
		Name:        name,
		NumSegments: numSegments,
		Segments:    &segments{b: make([]byte, numBytes), len: numSegments},
	}
}
