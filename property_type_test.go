// gosln.  An implementation of Semantic Link Network (SLN) in Go (Golang).
// Copyright (C) 2023  Yuan Gao
//
// This file is part of gosln.
//
// gosln is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package gosln_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/donyori/gosln"
)

type MyInt int

func TestPropertyTypeOf(t *testing.T) {
	intPtr := new(int)
	testCases := []struct {
		v    any
		want gosln.PropertyType
	}{
		{nil, 0},
		{false, gosln.Bool},
		{0, gosln.Int},
		{int8(0), gosln.Int8},
		{int16(0), gosln.Int16},
		{int32(0), gosln.Int32},
		{int64(0), gosln.Int64},
		{uint(0), gosln.Uint},
		{uint8(0), gosln.Uint8},
		{uint16(0), gosln.Uint16},
		{uint32(0), gosln.Uint32},
		{uint64(0), gosln.Uint64},
		{uintptr(0), gosln.Uintptr},
		{float32(0), gosln.Float32},
		{float64(0), gosln.Float64},
		{complex64(0), gosln.Complex64},
		{complex128(0), gosln.Complex128},
		{[]byte{}, gosln.Bytes},
		{"", gosln.String},
		{time.Time{}, gosln.Time},
		{MyInt(0), 0},
		{intPtr, 0},
		{gosln.Type{}, 0},
		{gosln.ID{}, 0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("vType=%T", tc.v), func(t *testing.T) {
			got := gosln.PropertyTypeOf(tc.v)
			if got != tc.want {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestPropertyType_Type(t *testing.T) {
	testCases := []struct {
		t     gosln.PropertyType
		wantV any
	}{
		{-1, nil},
		{0, nil},
		{gosln.Bool, false},
		{gosln.Int, 0},
		{gosln.Int8, int8(0)},
		{gosln.Int16, int16(0)},
		{gosln.Int32, int32(0)},
		{gosln.Int64, int64(0)},
		{gosln.Uint, uint(0)},
		{gosln.Uint8, uint8(0)},
		{gosln.Uint16, uint16(0)},
		{gosln.Uint32, uint32(0)},
		{gosln.Uint64, uint64(0)},
		{gosln.Uintptr, uintptr(0)},
		{gosln.Float32, float32(0)},
		{gosln.Float64, float64(0)},
		{gosln.Complex64, complex64(0)},
		{gosln.Complex128, complex128(0)},
		{gosln.Bytes, []byte{}},
		{gosln.String, ""},
		{gosln.Time, time.Time{}},
		{20, nil},
		{21, nil},
	}

	for _, tc := range testCases {
		var want reflect.Type
		if tc.wantV != nil {
			want = reflect.TypeOf(tc.wantV)
		}
		t.Run(fmt.Sprintf("i=%d", tc.t), func(t *testing.T) {
			got := tc.t.Type()
			if got != want {
				t.Errorf("got %v; want %v", got, want)
			}
		})
	}
}
