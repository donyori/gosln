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

	"github.com/donyori/gogo/errors"

	"github.com/donyori/gosln"
)

type MyInt int

func TestPropTypeOf(t *testing.T) {
	intPtr := new(int)
	testCases := []struct {
		v    any
		want gosln.PropType
	}{
		{nil, 0},
		{false, gosln.PTBool},
		{0, gosln.PTInt},
		{int8(0), gosln.PTInt8},
		{int16(0), gosln.PTInt16},
		{int32(0), gosln.PTInt32},
		{int64(0), gosln.PTInt64},
		{uint(0), gosln.PTUint},
		{uint8(0), gosln.PTUint8},
		{uint16(0), gosln.PTUint16},
		{uint32(0), gosln.PTUint32},
		{uint64(0), gosln.PTUint64},
		{uintptr(0), gosln.PTUintptr},
		{float32(0), gosln.PTFloat32},
		{float64(0), gosln.PTFloat64},
		{complex64(0), gosln.PTComplex64},
		{complex128(0), gosln.PTComplex128},
		{[]byte{}, gosln.PTBytes},
		{"", gosln.PTString},
		{time.Time{}, gosln.PTTime},
		{gosln.Date{}, gosln.PTDate},
		{MyInt(0), 0},
		{intPtr, 0},
		{gosln.Type{}, 0},
		{gosln.ID{}, 0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("vType=%T", tc.v), func(t *testing.T) {
			got := gosln.PropTypeOf(tc.v)
			if got != tc.want {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestPropType_GoType(t *testing.T) {
	testCases := []struct {
		t     gosln.PropType
		wantV any
	}{
		{-1, nil},
		{0, nil},
		{gosln.PTBool, false},
		{gosln.PTInt, 0},
		{gosln.PTInt8, int8(0)},
		{gosln.PTInt16, int16(0)},
		{gosln.PTInt32, int32(0)},
		{gosln.PTInt64, int64(0)},
		{gosln.PTUint, uint(0)},
		{gosln.PTUint8, uint8(0)},
		{gosln.PTUint16, uint16(0)},
		{gosln.PTUint32, uint32(0)},
		{gosln.PTUint64, uint64(0)},
		{gosln.PTUintptr, uintptr(0)},
		{gosln.PTFloat32, float32(0)},
		{gosln.PTFloat64, float64(0)},
		{gosln.PTComplex64, complex64(0)},
		{gosln.PTComplex128, complex128(0)},
		{gosln.PTBytes, []byte{}},
		{gosln.PTString, ""},
		{gosln.PTTime, time.Time{}},
		{gosln.PTDate, gosln.Date{}},
		{21, nil},
		{22, nil},
	}

	for _, tc := range testCases {
		var want reflect.Type
		if tc.wantV != nil {
			want = reflect.TypeOf(tc.wantV)
		}
		t.Run(fmt.Sprintf("i=%d", tc.t), func(t *testing.T) {
			got := tc.t.GoType()
			if got != want {
				t.Errorf("got %v; want %v", got, want)
			}
		})
	}
}

func TestPropTypeMap_Set(t *testing.T) {
	const (
		NoError int8 = iota
		InvalidPropName
		InvalidPropType
	)

	var pts []gosln.PropType
	pts = append(pts, 0, 1)
	for i := gosln.PropType(1); i.IsValid(); i++ {
		pts = append(pts, i+1)
	}

	testCases := make([]struct {
		propName         gosln.PropName
		propType         gosln.PropType
		wantPanicErrType int8
	}, 2*len(pts))
	var idx int
	for _, pn := range []gosln.PropName{gosln.MustNewPropName("prop"), {}} {
		for _, pt := range pts {
			testCases[idx].propName = pn
			testCases[idx].propType = pt
			if !pn.IsValid() {
				testCases[idx].wantPanicErrType = InvalidPropName
			} else if !pt.IsValid() {
				testCases[idx].wantPanicErrType = InvalidPropType
			}
			idx++
		}
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("propName=%v&propType=%v", tc.propName, tc.propType), func(t *testing.T) {
			ptm := gosln.NewPropTypeMap(1)
			defer func() {
				e := recover()
				if tc.wantPanicErrType == NoError {
					if e != nil {
						t.Error("panic -", e)
					}
					return
				} else if e == nil {
					t.Error("want panic but not")
					return
				}
				var target error
				err, ok := e.(error)
				switch {
				case !ok:
					t.Error("panic -", e)
					return
				case tc.wantPanicErrType == InvalidPropName:
					target = (*gosln.InvalidPropNameError)(nil)
				case tc.wantPanicErrType == InvalidPropType:
					target = (*gosln.InvalidPropTypeError)(nil)
				default:
					// This should never happen, but will act as a safeguard for later,
					// as a default value doesn't make sense here.
					t.Errorf("unknown wantPanicErrType %q", tc.wantPanicErrType)
					return
				}
				if !errors.As(err, &target) {
					t.Errorf("got error %v (%[1]T); want of type %T", err, target)
				}
			}()
			ptm.Set(tc.propName, tc.propType)
		})
	}
}
