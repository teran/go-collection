package mock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeNowMockString(t *testing.T) {
	r := require.New(t)

	ts := NewTimeNowMock()
	ts.On("Now").Return("2024-07-06T10:11:12Z").Once()

	r.Equal(time.Date(2024, 7, 6, 10, 11, 12, 0, time.UTC), ts.Now())
}

func TestTimeNowMockTime(t *testing.T) {
	r := require.New(t)

	ts := NewTimeNowMock()
	ts.On("Now").Return(time.Date(2024, 7, 6, 10, 11, 12, 0, time.UTC)).Once()

	r.Equal(time.Date(2024, 7, 6, 10, 11, 12, 0, time.UTC), ts.Now())
}
