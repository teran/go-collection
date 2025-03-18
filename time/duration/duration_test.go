package time

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationJSON(t *testing.T) {
	r := require.New(t)

	sample := `"15s"`
	var d Duration

	err := json.Unmarshal([]byte(sample), &d)
	r.NoError(err)
	r.Equal(15*Second, d)
	r.Equal(15*time.Second, d.TimeDuration())

	data, err := json.Marshal(d)
	r.NoError(err)
	r.Equal(`"15s"`, string(data))
}
