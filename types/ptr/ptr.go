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

import "time"

func Bool(b bool) *bool {
	return &b
}

func Duration(d time.Duration) *time.Duration {
	return &d
}

func Float32(f float32) *float32 {
	return &f
}

func Float64(f float64) *float64 {
	return &f
}

func Int(i int) *int {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func String(s string) *string {
	return &s
}

func Time(t time.Time) *time.Time {
	return &t
}

func Uint8(u uint8) *uint8 {
	return &u
}

func Uint16(u uint16) *uint16 {
	return &u
}

func Uint32(u uint32) *uint32 {
	return &u
}

func Uint64(u uint64) *uint64 {
	return &u
}

func Ptr[T any](in T) *T {
	return &in
}
