package choices

import (
	"math/rand"

	"github.com/Nordstrom/choices/util"
	"github.com/pkg/errors"
)

type Namespace struct {
	Name        string
	NumSegments int
	Segments    []byte
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
		Segments:    make([]byte, numBytes),
	}
}

func (n *Namespace) sampleSegments(num int) []byte {
	if num <= 0 {
		return make([]byte, len(n.Segments))
	}
	available := n.availableSegments()
	p := rand.Perm(len(available))
	out := make([]byte, len(n.Segments))
	if num > len(available) {
		num = len(available)
	}
	for i := 0; i < num; i++ {
		set(out, available[p[i]])
	}
	return out
}

func (n *Namespace) claimSegments(s []byte) error {
	if len(n.Segments) != len(s) {
		return errors.New("namespace and experiment have different number of segments")
	}
	for i, b := range s {
		n.Segments[i] |= b
	}
	return nil
}

func (n *Namespace) releaseSegments(s []byte) error {
	if len(n.Segments) != len(s) {
		return errors.New("namespace and experiment have different number of segments")
	}
	for i, b := range s {
		n.Segments[i] &= ^b
	}
	return nil
}

func (n *Namespace) availableSegments() []int {
	out := make([]int, 0, n.NumSegments)
	for i := range n.Segments {
		for shift := uint8(0); shift < 8; shift++ {
			if i*8+int(shift) > int(n.NumSegments) {
				break
			}
			if n.Segments[i]&(1<<shift) != 1<<shift {
				out = append(out, i*8+int(shift))
			}
		}
	}
	return out
}

func set(b []byte, index int) {
	i, pos := index/8, uint8(index%8)
	b[i] |= 1 << pos
}

func clear(b []byte, index int) {
	i, pos := index/8, uint8(index%8)
	b[i] &= ^(1 << pos)
}
