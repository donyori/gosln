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

func TestPropMapGet(t *testing.T) {
	name := gosln.MustNewPropName("gUpper")
	value := 'G'

	const (
		NoError int8 = iota
		PropNotExistError
		PropTypeError
	)

	pm := gosln.NewPropMap(1)
	err := gosln.PropMapSet(pm, name, value)
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
			fmt.Sprintf("type=rune&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[rune](pm, name)
			},
			value,
			NoError,
		},
		{
			fmt.Sprintf("type=int&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[int](pm, name)
			},
			int(value),
			NoError,
		},
		{
			fmt.Sprintf("type=int8&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[int8](pm, name)
			},
			int8(value),
			NoError,
		},
		{
			fmt.Sprintf("type=byte&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[byte](pm, name)
			},
			byte(value),
			NoError,
		},
		{
			fmt.Sprintf("type=uintptr&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[uintptr](pm, name)
			},
			uintptr(value),
			NoError,
		},
		{
			fmt.Sprintf("type=float32&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[float32](pm, name)
			},
			float32(value),
			NoError,
		},
		{
			fmt.Sprintf("type=string&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[string](pm, name)
			},
			string(value),
			NoError,
		},
		{
			fmt.Sprintf("type=rune&pm=<nil>&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[rune](nil, name)
			},
			nil,
			PropNotExistError,
		},
		{
			fmt.Sprintf("type=rune&pm=<empty>&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[rune](gosln.NewPropMap(0), name)
			},
			nil,
			PropNotExistError,
		},
		{
			fmt.Sprintf("type=rune&name=%+q", ""),
			func() (any, error) {
				return gosln.PropMapGet[rune](pm, gosln.PropName{})
			},
			nil,
			PropNotExistError,
		},
		{
			fmt.Sprintf("type=bool&name=%+q", ""),
			func() (any, error) {
				return gosln.PropMapGet[bool](pm, gosln.PropName{})
			},
			nil,
			PropNotExistError,
		},
		{
			fmt.Sprintf("type=bool&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[bool](pm, name)
			},
			nil,
			PropTypeError,
		},
		{
			fmt.Sprintf("type=complex128&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[complex128](pm, name)
			},
			nil,
			PropTypeError,
		},
		{
			fmt.Sprintf("type=[]byte&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[[]byte](pm, name)
			},
			nil,
			PropTypeError,
		},
		{
			fmt.Sprintf("type=time.Time&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[time.Time](pm, name)
			},
			nil,
			PropTypeError,
		},
		{
			fmt.Sprintf("type=gosln.Date&name=%+q", name),
			func() (any, error) {
				return gosln.PropMapGet[gosln.Date](pm, name)
			},
			nil,
			PropTypeError,
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
			case PropNotExistError:
				target = (*gosln.PropNotExistError)(nil)
			case PropTypeError:
				target = (*gosln.PropTypeError)(nil)
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

func TestPropMapGet_TimeAndDate(t *testing.T) {
	const Year int = 2023
	const Month = time.March
	const Day int = 12
	ti := time.Date(Year, Month, Day, 0, 0, 0, 0, time.UTC)
	date := gosln.DateOfYearMonthDay(Year, Month, Day)
	propName := gosln.MustNewPropName("date")

	t.Run("timeToDate", func(t *testing.T) {
		pm := gosln.NewPropMap(1)
		err := gosln.PropMapSet(pm, propName, ti)
		if err != nil {
			t.Fatal("set property -", err)
		}
		got, err := gosln.PropMapGet[gosln.Date](pm, propName)
		if err != nil {
			t.Errorf("got error (%v); want nil", err)
		} else if got != date {
			t.Errorf("got %v; want %v", got, date)
		}
	})

	t.Run("dateToTime", func(t *testing.T) {
		pm := gosln.NewPropMap(1)
		err := gosln.PropMapSet(pm, propName, date)
		if err != nil {
			t.Fatal("set property -", err)
		}
		got, err := gosln.PropMapGet[time.Time](pm, propName)
		if err != nil {
			t.Errorf("got error (%v); want nil", err)
		} else if !got.Equal(ti) {
			t.Errorf("got %v; want %v", got, ti)
		}
	})
}
