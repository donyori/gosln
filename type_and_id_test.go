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
	"strings"
	"testing"
	"time"

	"github.com/donyori/gosln"
)

func TestIsValidTypeString(t *testing.T) {
	longestName := strings.Repeat("A", 65535)
	testCases := []struct {
		t    string
		want bool
	}{
		{"A", true},
		{"ABC", true},
		{"Ab_4", true},
		{"P0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", true},
		{longestName, true},
		{"", false},
		{"a", false},
		{"aBC", false},
		{"0BC", false},
		{"_BC", false},
		{"-BC", false},
		{"AB-", false},
		{"AB" + string('0'-1), false},
		{"AB" + string('9'+1), false},
		{"AB" + string('A'-1), false},
		{"AB" + string('Z'+1), false},
		{"AB" + string('a'-1), false},
		{"AB" + string('z'+1), false},
		{"AB-D", false},
		{"SLN", false},
		{"SLNType", false},
		{longestName + "A", false},
	}

	for _, tc := range testCases {
		name := tc.t
		if len(name) > 40 {
			name = fmt.Sprintf("%s...%s(len=%d)", name[:8], name[len(name)-8:], len(name))
		}
		t.Run(fmt.Sprintf("type=%+q", name), func(t *testing.T) {
			got := gosln.IsValidTypeString(tc.t)
			if got != tc.want {
				t.Errorf("got %t; want %t", got, tc.want)
			}
		})
	}
}

func TestNewID(t *testing.T) {
	date := gosln.DateOfYearMonthDay(2023, time.March, 12)
	typ1 := gosln.MustNewType("TestType_1")
	typ2 := gosln.MustNewType("TestType_2")
	testCases := []struct {
		t         gosln.Type
		i         int64
		wantStr   string
		wantPanic bool
	}{
		{gosln.Type{}, 0, "", false},
		{gosln.Type{}, 1, "", false},
		{typ1, 0, "TestType_1#2023-071-0", false},
		{typ1, 1, "TestType_1#2023-071-1", false},
		{typ1, 9, "TestType_1#2023-071-9", false},
		{typ1, 10, "TestType_1#2023-071-A", false},
		{typ1, 35, "TestType_1#2023-071-Z", false},
		{typ1, 36, "TestType_1#2023-071-a", false},
		{typ1, 61, "TestType_1#2023-071-z", false},
		{typ1, 62, "TestType_1#2023-071--", false},
		{typ1, 63, "TestType_1#2023-071-_", false},
		{typ1, 64, "TestType_1#2023-071-00", false},
		{typ1, 65, "TestType_1#2023-071-10", false},
		{typ1, 73, "TestType_1#2023-071-90", false},
		{typ1, 74, "TestType_1#2023-071-A0", false},
		{typ1, 99, "TestType_1#2023-071-Z0", false},
		{typ1, 100, "TestType_1#2023-071-a0", false},
		{typ1, 125, "TestType_1#2023-071-z0", false},
		{typ1, 126, "TestType_1#2023-071--0", false},
		{typ1, 127, "TestType_1#2023-071-_0", false},
		{typ1, 128, "TestType_1#2023-071-01", false},
		{typ1, 129, "TestType_1#2023-071-11", false},
		{typ1, 191, "TestType_1#2023-071-_1", false},
		{typ1, 192, "TestType_1#2023-071-02", false},
		{typ1, 193, "TestType_1#2023-071-12", false},
		{typ1, 255, "TestType_1#2023-071-_2", false},
		{typ1, 256, "TestType_1#2023-071-03", false},
		{typ1, 640, "TestType_1#2023-071-09", false},
		{typ1, 704, "TestType_1#2023-071-0A", false},
		{typ1, 2304, "TestType_1#2023-071-0Z", false},
		{typ1, 2368, "TestType_1#2023-071-0a", false},
		{typ1, 3968, "TestType_1#2023-071-0z", false},
		{typ1, 4032, "TestType_1#2023-071-0-", false},
		{typ1, 4096, "TestType_1#2023-071-0_", false},
		{typ1, 4159, "TestType_1#2023-071-__", false},
		{typ1, 4160, "TestType_1#2023-071-000", false},
		{typ1, 4161, "TestType_1#2023-071-100", false},
		{typ1, 8256, "TestType_1#2023-071-001", false},
		{typ1, 262208, "TestType_1#2023-071-00_", false},
		{typ1, 266304, "TestType_1#2023-071-0000", false},
		{typ2, 0, "TestType_2#2023-071-0", false},
		{typ2, 1, "TestType_2#2023-071-1", false},
		{gosln.Type{}, -1, "", true},
		{typ1, -1, "", true},
		{typ2, -1, "", true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("type=%+q&i=%d", tc.t, tc.i), func(t *testing.T) {
			defer func() {
				e := recover()
				if tc.wantPanic {
					if e == nil {
						t.Error("want panic but not")
					}
				} else if e != nil {
					t.Error("panic -", e)
				}
			}()
			id := gosln.NewID(tc.t, date, tc.i)
			if id.String() != tc.wantStr {
				t.Errorf("got %v; want %s", id, tc.wantStr)
			}
			if isValid, wantValid := id.IsValid(), tc.wantStr != ""; isValid != wantValid {
				t.Errorf("got IsValid %t; want %t", isValid, wantValid)
			}
			if typ := id.Type(); typ != tc.t {
				t.Errorf("got Type %v; want %v", typ, tc.t)
			}
		})
	}
}
