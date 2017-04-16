package rand

import (
	"math/rand"
	"sync"
	"time"
)

var rng = struct {
	sync.Mutex
	rand *rand.Rand
}{
	rand: rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
}

// Intn generates an integer in range [0,max).
// By design this should panic if input is invalid, <= 0.
func Intn(max int) int {
	rng.Lock()
	defer rng.Unlock()
	return rng.rand.Intn(max)
}

var alphanums = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// String generates a string with length n out of alpha-numeric characters.
func String(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = alphanums[Intn(len(alphanums))]
	}
	return string(b)
}
