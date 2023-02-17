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
	"testing"
	"time"

	"github.com/donyori/gogo/errors"

	"github.com/donyori/gosln"
)

func TestGetProperty(t *testing.T) {
	const Name = "gUpper"
	value := 'G'

	const (
		NoError int8 = iota
		PropertyNotExistError
		PropertyTypeError
	)

	var pm gosln.PropertyMap
	err := gosln.SetProperty(&pm, Name, value)
	if err != nil {
		t.Fatal("set property -", err)
	}

	testCases := []struct {
		name        string
		f           func() (any, error)
		wantV       any
		wantErrType int8
	}{
		{
			fmt.Sprintf("type=rune&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[rune](&pm, Name)
			},
			value,
			NoError,
		},
		{
			fmt.Sprintf("type=int&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[int](&pm, Name)
			},
			int(value),
			NoError,
		},
		{
			fmt.Sprintf("type=int8&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[int8](&pm, Name)
			},
			int8(value),
			NoError,
		},
		{
			fmt.Sprintf("type=byte&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[byte](&pm, Name)
			},
			byte(value),
			NoError,
		},
		{
			fmt.Sprintf("type=uintptr&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[uintptr](&pm, Name)
			},
			uintptr(value),
			NoError,
		},
		{
			fmt.Sprintf("type=float32&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[float32](&pm, Name)
			},
			float32(value),
			NoError,
		},
		{
			fmt.Sprintf("type=string&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[string](&pm, Name)
			},
			string(value),
			NoError,
		},
		{
			fmt.Sprintf("type=rune&name=%+q", ""),
			func() (any, error) {
				return gosln.GetProperty[rune](&pm, "")
			},
			nil,
			PropertyNotExistError,
		},
		{
			fmt.Sprintf("type=bool&name=%+q", ""),
			func() (any, error) {
				return gosln.GetProperty[bool](&pm, "")
			},
			nil,
			PropertyNotExistError,
		},
		{
			fmt.Sprintf("type=bool&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[bool](&pm, Name)
			},
			nil,
			PropertyTypeError,
		},
		{
			fmt.Sprintf("type=complex128&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[complex128](&pm, Name)
			},
			nil,
			PropertyTypeError,
		},
		{
			fmt.Sprintf("type=[]byte&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[[]byte](&pm, Name)
			},
			nil,
			PropertyTypeError,
		},
		{
			fmt.Sprintf("type=time.Time&name=%+q", Name),
			func() (any, error) {
				return gosln.GetProperty[time.Time](&pm, Name)
			},
			nil,
			PropertyTypeError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := tc.f()
			var target error
			switch tc.wantErrType {
			case NoError:
				if err != nil {
					t.Errorf("got error (%v); want nil", err)
				} else if v != tc.wantV {
					t.Errorf("got value %v (%[1]T); want %v (%[2]T)", v, tc.wantV)
				}
				return
			case PropertyNotExistError:
				target = new(gosln.PropertyNotExistError)
			case PropertyTypeError:
				target = new(gosln.PropertyTypeError)
			default:
				// This should never happen, but will act as a safeguard for later,
				// as a default value doesn't make sense here.
				t.Fatalf("unknown wantErrType %q", tc.wantErrType)
			}
			if !errors.As(err, &target) {
				t.Errorf("got error %v (%[1]T); want of type %T", err, target)
			}
		})
	}
}
