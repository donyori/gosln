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

package gosln

import (
	"fmt"
	"time"
)

// Date represents a date (an instant in time with day precision).
//
// It records the year (of the Common Era (CE)) and the day within the year
// in Universal Coordinated Time (UTC).
type Date struct {
	year    int
	yearDay int16
}

// NowDate returns the current date (in UTC instead of the local time zone).
func NowDate() Date {
	return DateOf(time.Now())
}

// DateOf returns the date of the specified time t.
func DateOf(t time.Time) Date {
	t = t.UTC()
	return Date{
		year:    t.Year(),
		yearDay: int16(t.YearDay()),
	}
}

// IsZero reports whether the date is a zero-value Date.
func (d Date) IsZero() bool {
	return d.year == 0 && d.yearDay == 0
}

// GoTime returns the time.Time corresponding to the date,
// whose hour, minute, second, and nanosecond are 0,
// and the location is UTC.
func (d Date) GoTime() time.Time {
	return time.Date(d.year, 0, int(d.yearDay), 0, 0, 0, 0, time.UTC)
}

// Year returns the year of the date.
func (d Date) Year() int {
	return d.year
}

// YearDay returns the day of the year specified by the date,
// in the range [1,365] for non-leap years, and [1,366] in leap years.
func (d Date) YearDay() int {
	return int(d.yearDay)
}

// Equal reports whether this date is the same as the specified date.
func (d Date) Equal(date Date) bool {
	return d.year == date.year && d.yearDay == date.yearDay
}

// Before reports whether this date is before the specified date.
func (d Date) Before(date Date) bool {
	return d.year < date.year ||
		d.year == date.year && d.yearDay < date.yearDay
}

// After reports whether this date is after the specified date.
func (d Date) After(date Date) bool {
	return d.year > date.year ||
		d.year == date.year && d.yearDay > date.yearDay
}

// Compare compares this date (denoted by x)
// and the specified date (denoted by y).
//
// If x is before y, it returns -1;
// if x is after y, it returns +1;
// if x and y are the same, it returns 0.
func (d Date) Compare(date Date) int {
	switch {
	case d.year < date.year:
		return -1
	case d.year > date.year:
		return 1
	case d.yearDay < date.yearDay:
		return -1
	case d.yearDay > date.yearDay:
		return 1
	}
	return 0
}

// String formats the date in the form of
//
//	<YEAR> "-" <YEAR-DAY>
//
// where <YEAR> is a decimal integer with no padding,
// and <YEAR-DAY> is a 3-digit decimal integer padding with "0".
//
// The result is the same as fmt.Sprintf("%d-%03d", d.Year(), d.YearDay()).
func (d Date) String() string {
	return fmt.Sprintf("%d-%03d", d.year, d.yearDay)
}
