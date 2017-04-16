package util

import (
	"fmt"
	"github.com/Nordstrom/choices/util/rand"
)

type nameGenerator struct {
	maxLength int
	randomLength int
}

func (n *nameGenerator) GenerateName(base string) string {
	maxGeneratedNameLength := n.maxLength - n.randomLength
	if len(base) > maxGeneratedNameLength-1 {
		base = base[:maxGeneratedNameLength-1]
	}
	return fmt.Sprintf("%s%s", base, rand.String(n.randomLength))
}

var BasicNameGenerator = nameGenerator{maxLength: 64, randomLength: 5}
