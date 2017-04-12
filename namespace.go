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

func (n *Namespace) ToNamespace() *storage.Namespace {
	return &storage.Namespace{
		Name:        n.Name,
		NumSegments: int64(n.NumSegments),
		Segments:    n.Segments.ToSegments(),
	}
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
