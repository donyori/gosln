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
	year, yearDay int
}

// NowDate returns the current date (in UTC instead of the local time zone).
func NowDate() Date {
	return DateOf(time.Now())
}

// DateOf returns the date specified by the time t, converted to UTC.
func DateOf(t time.Time) Date {
	t = t.UTC()
	return Date{
		year:    t.Year(),
		yearDay: t.YearDay(),
	}
}

// DateOfYearMonthDay returns the date specified by the year, month, and day
// (in UTC instead of the local time zone).
//
// Similar to the function time.Date,
// month and day may be outside their usual ranges.
// DateOfYearMonthDay normalizes them during the conversion.
// For example, October 32 converts to November 1.
func DateOfYearMonthDay(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return Date{
		year:    t.Year(),
		yearDay: t.YearDay(),
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
	// Set the month to January (1) rather than 0;
	// 0 is normalized to December last year.
	return time.Date(d.year, time.January, d.yearDay, 0, 0, 0, 0, time.UTC)
}

// Year returns the year of the date.
func (d Date) Year() int {
	return d.year
}

// Month returns the month of the year specified by the date.
func (d Date) Month() time.Month {
	return d.GoTime().Month()
}

// Day returns the day of the month specified by the date.
func (d Date) Day() int {
	return d.GoTime().Day()
}

// YearDay returns the day of the year specified by the date,
// in the range [1,365] for non-leap years, and [1,366] in leap years.
func (d Date) YearDay() int {
	return d.yearDay
}

// Weekday returns the day of the week specified by the date.
func (d Date) Weekday() time.Weekday {
	return d.GoTime().Weekday()
}

// YearMonthDay returns the year, month, and day specified by the date.
func (d Date) YearMonthDay() (year int, month time.Month, day int) {
	return d.GoTime().Date()
}

// ISOWeek returns the ISO 8601 year and week number specified by the date.
//
// Week ranges from 1 to 53.
// Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1,
// and Dec 29 to Dec 31 might belong to week 1 of year n+1.
func (d Date) ISOWeek() (year int, week int) {
	return d.GoTime().ISOWeek()
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
	a, b := d.year, date.year
	if a == b {
		a, b = d.yearDay, date.yearDay
	}
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// Add returns the date after the specified duration since this date.
func (d Date) Add(duration time.Duration) Date {
	t := d.GoTime().Add(duration)
	return Date{
		year:    t.Year(),
		yearDay: t.YearDay(),
	}
}

// AddYearMonthDay returns the date corresponding to adding
// the specified number of years, months, and days to this date.
func (d Date) AddYearMonthDay(years, months, days int) Date {
	t := time.Date(
		d.year+years, time.January+time.Month(months), d.yearDay+days,
		0, 0, 0, 0, time.UTC,
	)
	return Date{
		year:    t.Year(),
		yearDay: t.YearDay(),
	}
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
