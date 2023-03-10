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

func TestIsValidPropNameString(t *testing.T) {
	longestName := strings.Repeat("a", 65535)
	testCases := []struct {
		propName string
		want     bool
	}{
		{"a", true},
		{"abc", true},
		{"aB_4", true},
		{"p0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", true},
		{longestName, true},
		{"", false},
		{"A", false},
		{"Abc", false},
		{"0bc", false},
		{"_bc", false},
		{"-bc", false},
		{"ab-", false},
		{"ab" + string('0'-1), false},
		{"ab" + string('9'+1), false},
		{"ab" + string('A'-1), false},
		{"ab" + string('Z'+1), false},
		{"ab" + string('a'-1), false},
		{"ab" + string('z'+1), false},
		{"ab-d", false},
		{"sln", false},
		{"slnID", false},
		{"slnType", false},
		{longestName + "a", false},
	}

	for _, tc := range testCases {
		name := tc.propName
		if len(name) > 40 {
			name = fmt.Sprintf("%s...%s(len=%d)", name[:8], name[len(name)-8:], len(name))
		}
		t.Run(fmt.Sprintf("name=%+q", name), func(t *testing.T) {
			got := gosln.IsValidPropNameString(tc.propName)
			if got != tc.want {
				t.Errorf("got %t; want %t", got, tc.want)
			}
		})
	}
}
