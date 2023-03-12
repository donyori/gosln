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

	"github.com/donyori/gosln"
)

func TestNowDate(t *testing.T) {
	var now time.Time
	var nowDate gosln.Date
	for {
		now = time.Now()
		nowDate = gosln.NowDate()
		// Check time.Now() again to ensure that the execution of
		// the above statements did not cross days.
		u := time.Now()
		if now.Year() == u.Year() && now.YearDay() == u.YearDay() {
			break
		}
	}
	now = now.UTC()
	gotYear, gotYearDay := nowDate.Year(), nowDate.YearDay()
	wantYear, wantYearDay := now.Year(), now.YearDay()
	if gotYear != wantYear || gotYearDay != wantYearDay {
		t.Errorf("got Year %d, YearDay %d; want Year %d, YearDay %d",
			gotYear, gotYearDay, wantYear, wantYearDay)
	}
}

func TestDateOfAndGoTime(t *testing.T) {
	cst := time.FixedZone("CST", 8*60*60)
	times := []time.Time{
		{},
		time.Date(1, time.January, 1, 0, 0, 0, 0, cst),
		time.Unix(0, 0).UTC(),
		time.Unix(0, 0).In(cst).Add(time.Hour * -8),
		time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
		time.Date(0, 0, 0, 0, 0, 0, 0, cst),
		time.Date(2023, time.March, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.March, 12, 0, 0, 0, 0, cst),
		time.Date(2023, time.January, 365, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.January, 365, 0, 0, 0, 0, cst),
		time.Date(2020, time.February, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.February, 29, 0, 0, 0, 0, cst),
	}

	for _, x := range times {
		t.Run(fmt.Sprintf("time=%v", x), func(t *testing.T) {
			cp := x
			date := gosln.DateOf(cp)
			got := date.GoTime()
			want := cp.UTC()
			gotYear, gotYearDay := got.Year(), got.YearDay()
			wantYear, wantYearDay := want.Year(), want.YearDay()
			if gotYear != wantYear || gotYearDay != wantYearDay {
				t.Errorf("got Year %d, YearDay %d; want Year %d, YearDay %d",
					gotYear, gotYearDay, wantYear, wantYearDay)
			}
		})
	}
}

func TestDateOfYearMonthDay(t *testing.T) {
	testCases := []struct {
		year        int
		month       time.Month
		day         int
		wantYear    int
		wantYearDay int
	}{
		{1, time.January, 1, 1, 1},
		{0, 0, 0, -1, 334},
		{2023, time.March, 12, 2023, 71},
		{2023, time.January, 71, 2023, 71},
		{2022, time.February, 405, 2023, 71},
		{2023, time.December, 31, 2023, 365},
		{2023, 13, 0, 2023, 365},
		{2023, time.December, 32, 2024, 1},
		{2020, time.December, 31, 2020, 366},
		{2020, 13, 0, 2020, 366},
		{2020, time.December, 32, 2021, 1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("year=%d&month=%v&day=%d", tc.year, tc.month, tc.day), func(t *testing.T) {
			date := gosln.DateOfYearMonthDay(tc.year, tc.month, tc.day)
			gotYear, gotYearDay := date.Year(), date.YearDay()
			if gotYear != tc.wantYear || gotYearDay != tc.wantYearDay {
				t.Errorf("got Year %d, YearDay %d; want Year %d, YearDay %d",
					gotYear, gotYearDay, tc.wantYear, tc.wantYearDay)
			}
		})
	}
}
