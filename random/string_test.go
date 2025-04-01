package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	r := require.New(t)

	samples := map[string]struct{}{}
	for i := 0; i < 1024; i++ {
		samples[String(AlphaNumericWithSpecial, 128)] = struct{}{}
	}

	r.Len(samples, 1024)
}
