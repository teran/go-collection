package ptr

// Copyright 2023 Igor Shishkin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func (s *ptrTestSuite) TestBool() {
	b := Bool(true)
	s.Require().Equal(true, *b)
}

func (s *ptrTestSuite) TestDuration() {
	d := Duration(1 * time.Minute)
	s.Require().Equal(1*time.Minute, *d)
}

func (s *ptrTestSuite) TestFloat32() {
	f := Float32(123)
	s.Require().Equal(float32(123), *f)
}

func (s *ptrTestSuite) TestFloat64() {
	f := Float64(123)
	s.Require().Equal(float64(123), *f)
}

func (s *ptrTestSuite) TestInt() {
	i := Int(123)
	s.Require().Equal(int(123), *i)
}

func (s *ptrTestSuite) TestInt64() {
	i := Int64(123)
	s.Require().Equal(int64(123), *i)
}

func (s *ptrTestSuite) TestString() {
	st := String("test")
	s.Require().Equal("test", *st)
}

func (s *ptrTestSuite) TestTime() {
	now := time.Now()
	t := Time(now)
	s.Require().Equal(now, *t)
}

func (s *ptrTestSuite) TestUint8() {
	u := Uint8(123)
	s.Require().Equal(uint8(123), *u)
}

func (s *ptrTestSuite) TestUint16() {
	u := Uint16(123)
	s.Require().Equal(uint16(123), *u)
}

func (s *ptrTestSuite) TestUint32() {
	u := Uint32(123)
	s.Require().Equal(uint32(123), *u)
}

func (s *ptrTestSuite) TestUint64() {
	u := Uint64(123)
	s.Require().Equal(uint64(123), *u)
}

func (s *ptrTestSuite) TestPtr() {
	vi := Ptr(123)
	s.Require().Equal(int(123), *vi)

	vf := Ptr(123.321)
	s.Require().Equal(float64(123.321), *vf)

	vs := Ptr("test string")
	s.Require().Equal("test string", *vs)
}

// ========================================================================
// Test suite setup
// ========================================================================
type ptrTestSuite struct {
	suite.Suite
}

func TestPtrTestSuite(t *testing.T) {
	suite.Run(t, &ptrTestSuite{})
}
