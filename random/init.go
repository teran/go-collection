package random

import (
	"math/rand"
	"time"
)

var rng *rand.Rand

type Random interface {
	Int63() int64
	Uint32() uint32
	Int31() int32
	Int() int
	Int63n(n int64) int64
	Int31n(n int32) int32
	Intn(n int) int
	Float64() float64
	Float32() float32
	Perm(n int) []int
	Read([]byte) (int, error)
}

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // non-security random strings
}

func GetRand() Random {
	return rng
}
