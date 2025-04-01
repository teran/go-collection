package random

import "github.com/stretchr/testify/mock"

var _ Random = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Int63() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *Mock) Uint32() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *Mock) Int31() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *Mock) Int() int {
	args := m.Called()
	return args.Int(0)
}

func (m *Mock) Int63n(n int64) int64 {
	args := m.Called(n)
	return args.Get(0).(int64)
}

func (m *Mock) Int31n(n int32) int32 {
	args := m.Called(n)
	return args.Get(0).(int32)
}

func (m *Mock) Intn(n int) int {
	args := m.Called(n)
	return args.Int(0)
}

func (m *Mock) Float64() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *Mock) Float32() float32 {
	args := m.Called()
	return args.Get(0).(float32)
}

func (m *Mock) Perm(n int) []int {
	args := m.Called(n)
	return args.Get(0).([]int)
}

func (m *Mock) Read(buf []byte) (int, error) {
	args := m.Called()
	outBuf := args.Get(0).([]byte)
	return copy(buf, outBuf), args.Error(2)
}
