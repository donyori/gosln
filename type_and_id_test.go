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
	nowDate := gosln.NowDate()
	typ1 := gosln.MustNewType("TestType_1")
	typ2 := gosln.MustNewType("TestType_2")
	testCases := []struct {
		t             gosln.Type
		i             int64
		wantStrLayout string
		wantPanic     bool
	}{
		{gosln.Type{}, 0, "", false},
		{gosln.Type{}, 1, "", false},
		{typ1, 0, "TestType_1#%v-0", false},
		{typ1, 1, "TestType_1#%v-1", false},
		{typ1, 9, "TestType_1#%v-9", false},
		{typ1, 10, "TestType_1#%v-A", false},
		{typ1, 35, "TestType_1#%v-Z", false},
		{typ1, 36, "TestType_1#%v-a", false},
		{typ1, 61, "TestType_1#%v-z", false},
		{typ1, 62, "TestType_1#%v--", false},
		{typ1, 63, "TestType_1#%v-_", false},
		{typ1, 64, "TestType_1#%v-00", false},
		{typ1, 65, "TestType_1#%v-10", false},
		{typ1, 73, "TestType_1#%v-90", false},
		{typ1, 74, "TestType_1#%v-A0", false},
		{typ1, 99, "TestType_1#%v-Z0", false},
		{typ1, 100, "TestType_1#%v-a0", false},
		{typ1, 125, "TestType_1#%v-z0", false},
		{typ1, 126, "TestType_1#%v--0", false},
		{typ1, 127, "TestType_1#%v-_0", false},
		{typ1, 128, "TestType_1#%v-01", false},
		{typ1, 129, "TestType_1#%v-11", false},
		{typ1, 191, "TestType_1#%v-_1", false},
		{typ1, 192, "TestType_1#%v-02", false},
		{typ1, 193, "TestType_1#%v-12", false},
		{typ1, 255, "TestType_1#%v-_2", false},
		{typ1, 256, "TestType_1#%v-03", false},
		{typ1, 640, "TestType_1#%v-09", false},
		{typ1, 704, "TestType_1#%v-0A", false},
		{typ1, 2304, "TestType_1#%v-0Z", false},
		{typ1, 2368, "TestType_1#%v-0a", false},
		{typ1, 3968, "TestType_1#%v-0z", false},
		{typ1, 4032, "TestType_1#%v-0-", false},
		{typ1, 4096, "TestType_1#%v-0_", false},
		{typ1, 4159, "TestType_1#%v-__", false},
		{typ1, 4160, "TestType_1#%v-000", false},
		{typ1, 4161, "TestType_1#%v-100", false},
		{typ1, 8256, "TestType_1#%v-001", false},
		{typ1, 262208, "TestType_1#%v-00_", false},
		{typ1, 266304, "TestType_1#%v-0000", false},
		{typ2, 0, "TestType_2#%v-0", false},
		{typ2, 1, "TestType_2#%v-1", false},
		{gosln.Type{}, -1, "", true},
		{typ1, -1, "", true},
		{typ2, -1, "", true},
	}

	for _, tc := range testCases {
		wantStr := tc.wantStrLayout
		if wantStr != "" {
			wantStr = fmt.Sprintf(wantStr, nowDate)
		}
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
			id := gosln.NewID(tc.t, nowDate, tc.i)
			if id.String() != wantStr {
				t.Errorf("got %v; want %s", id, wantStr)
			}
			if isValid, wantValid := id.IsValid(), wantStr != ""; isValid != wantValid {
				t.Errorf("got IsValid %t; want %t", isValid, wantValid)
			}
			if typ := id.Type(); typ != tc.t {
				t.Errorf("got Type %v; want %v", typ, tc.t)
			}
		})
	}
}
