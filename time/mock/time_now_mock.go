package time

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

type TimeNowMock struct {
	mock.Mock
}

func NewTimeNowMock() *TimeNowMock {
	return &TimeNowMock{}
}

func (m *TimeNowMock) Now() time.Time {
	args := m.Called()

	switch v := args.Get(0).(type) {
	case time.Time:
		return v
	case string:
		ts, err := time.Parse(time.RFC3339, v)
		if err != nil {
			panic(err)
		}

		return ts
	default:
		panic(fmt.Sprintf("unexpected timestamp time %T", v))
	}
}
