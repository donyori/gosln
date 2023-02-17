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

	"github.com/donyori/gosln"
)

func TestNewID(t *testing.T) {
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
		{typ1, 0, "TestType_1#0", false},
		{typ1, 1, "TestType_1#1", false},
		{typ1, 9, "TestType_1#9", false},
		{typ1, 10, "TestType_1#A", false},
		{typ1, 35, "TestType_1#Z", false},
		{typ1, 36, "TestType_1#a", false},
		{typ1, 61, "TestType_1#z", false},
		{typ1, 62, "TestType_1#-", false},
		{typ1, 63, "TestType_1#_", false},
		{typ1, 64, "TestType_1#00", false},
		{typ1, 65, "TestType_1#10", false},
		{typ1, 73, "TestType_1#90", false},
		{typ1, 74, "TestType_1#A0", false},
		{typ1, 99, "TestType_1#Z0", false},
		{typ1, 100, "TestType_1#a0", false},
		{typ1, 125, "TestType_1#z0", false},
		{typ1, 126, "TestType_1#-0", false},
		{typ1, 127, "TestType_1#_0", false},
		{typ1, 128, "TestType_1#01", false},
		{typ1, 129, "TestType_1#11", false},
		{typ1, 191, "TestType_1#_1", false},
		{typ1, 192, "TestType_1#02", false},
		{typ1, 193, "TestType_1#12", false},
		{typ1, 255, "TestType_1#_2", false},
		{typ1, 256, "TestType_1#03", false},
		{typ1, 640, "TestType_1#09", false},
		{typ1, 704, "TestType_1#0A", false},
		{typ1, 2304, "TestType_1#0Z", false},
		{typ1, 2368, "TestType_1#0a", false},
		{typ1, 3968, "TestType_1#0z", false},
		{typ1, 4032, "TestType_1#0-", false},
		{typ1, 4096, "TestType_1#0_", false},
		{typ1, 4159, "TestType_1#__", false},
		{typ1, 4160, "TestType_1#000", false},
		{typ1, 4161, "TestType_1#100", false},
		{typ1, 8256, "TestType_1#001", false},
		{typ1, 262208, "TestType_1#00_", false},
		{typ1, 266304, "TestType_1#0000", false},
		{typ2, 0, "TestType_2#0", false},
		{typ2, 1, "TestType_2#1", false},
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
			id := gosln.NewID(tc.t, tc.i)
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
