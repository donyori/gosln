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
	"reflect"
	"time"
)

// PropertyType represents the type of property.
type PropertyType int8

const (
	Bool            PropertyType = 1 + iota // bool
	Int                                     // int
	Int8                                    // int8
	Int16                                   // int16
	Int32                                   // int32
	Int64                                   // int64
	Uint                                    // uint
	Uint8                                   // uint8
	Uint16                                  // uint16
	Uint32                                  // uint32
	Uint64                                  // uint64
	Uintptr                                 // uintptr
	Float32                                 // float32
	Float64                                 // float64
	Complex64                               // complex64
	Complex128                              // complex128
	Bytes                                   // []byte
	String                                  // string
	Time                                    // time.Time
	maxPropertyType                         // PropertyType(20)
)

// Before running the following command, please make sure the numeric value
// in the line comment of maxPropertyType is correct.
//
//go:generate stringer -type=PropertyType -output=property_type_string.go -linecomment

// PropertyTypeOf returns the property type of the value v.
//
// It returns 0 if v does not conform to PropertyValue.
func PropertyTypeOf(v any) PropertyType {
	return propertyTypeOfMap[reflect.TypeOf(v)]
}

// IsValid reports whether the property type is known.
func (i PropertyType) IsValid() bool {
	return i > 0 && i < maxPropertyType
}

// Type returns the reflect.Type corresponding to the property type.
//
// It returns nil if the property type is invalid.
func (i PropertyType) Type() reflect.Type {
	if i > 0 && i < maxPropertyType {
		return propertyTypes[i-1]
	}
	return nil
}

// IsConvertibleTo reports whether the property type i can convert to type t.
func (i PropertyType) IsConvertibleTo(t PropertyType) bool {
	if i <= 0 || i >= maxPropertyType || t <= 0 || t >= maxPropertyType {
		return false
	}
	return propertyTypes[i-1].ConvertibleTo(propertyTypes[t-1])
}

// IsSignedInteger reports whether the property type is a signed integer,
// including int, int8, int16, int32 (rune), and int64.
func (i PropertyType) IsSignedInteger() bool {
	switch i {
	case Int, Int8, Int16, Int32, Int64:
		return true
	}
	return false
}

// IsUnsignedInteger reports whether the property type is an unsigned integer,
// including uint, uint8 (byte), uint16, uint32, uint64, and uintptr.
func (i PropertyType) IsUnsignedInteger() bool {
	switch i {
	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		return true
	}
	return false
}

// IsInteger reports whether the property type is an integer,
// including int, int8, int16, int32 (rune), int64,
// uint, uint8 (byte), uint16, uint32, uint64, and uintptr.
func (i PropertyType) IsInteger() bool {
	switch i {
	case Int, Int8, Int16, Int32, Int64,
		Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		return true
	}
	return false
}

// IsFloat reports whether the property type is a floating number,
// including float32 and float64.
func (i PropertyType) IsFloat() bool {
	switch i {
	case Float32, Float64:
		return true
	}
	return false
}

// IsRealNumber reports whether the property type is a real number,
// including int, int8, int16, int32 (rune), int64,
// uint, uint8 (byte), uint16, uint32, uint64, uintptr,
// float32, and float64.
func (i PropertyType) IsRealNumber() bool {
	switch i {
	case Int, Int8, Int16, Int32, Int64,
		Uint, Uint8, Uint16, Uint32, Uint64, Uintptr,
		Float32, Float64:
		return true
	}
	return false
}

// IsComplex reports whether the property type is a complex number,
// including complex64 and complex128.
func (i PropertyType) IsComplex() bool {
	switch i {
	case Complex64, Complex128:
		return true
	}
	return false
}

// IsNumeric reports whether the property type is a number,
// including int, int8, int16, int32 (rune), int64,
// uint, uint8 (byte), uint16, uint32, uint64, uintptr,
// float32, float64, complex64, and complex128.
func (i PropertyType) IsNumeric() bool {
	switch i {
	case Int, Int8, Int16, Int32, Int64,
		Uint, Uint8, Uint16, Uint32, Uint64, Uintptr,
		Float32, Float64,
		Complex64, Complex128:
		return true
	}
	return false
}

// IsByteString reports whether the property type is a byte string,
// including []byte and string.
func (i PropertyType) IsByteString() bool {
	switch i {
	case Bytes, String:
		return true
	}
	return false
}

var (
	// propertyTypes is an array of reflect.Type corresponding to PropertyType.
	propertyTypes [maxPropertyType - 1]reflect.Type
	// propertyTypeOfMap is a map from reflect.Type to PropertyType,
	//used by PropertyTypeOf.
	propertyTypeOfMap map[reflect.Type]PropertyType
)

func init() {
	propertyTypes[Bool-1] = reflect.TypeOf(false)
	propertyTypes[Int-1] = reflect.TypeOf(0)
	propertyTypes[Int8-1] = reflect.TypeOf(int8(0))
	propertyTypes[Int16-1] = reflect.TypeOf(int16(0))
	propertyTypes[Int32-1] = reflect.TypeOf(int32(0))
	propertyTypes[Int64-1] = reflect.TypeOf(int64(0))
	propertyTypes[Uint-1] = reflect.TypeOf(uint(0))
	propertyTypes[Uint8-1] = reflect.TypeOf(uint8(0))
	propertyTypes[Uint16-1] = reflect.TypeOf(uint16(0))
	propertyTypes[Uint32-1] = reflect.TypeOf(uint32(0))
	propertyTypes[Uint64-1] = reflect.TypeOf(uint64(0))
	propertyTypes[Uintptr-1] = reflect.TypeOf(uintptr(0))
	propertyTypes[Float32-1] = reflect.TypeOf(float32(0))
	propertyTypes[Float64-1] = reflect.TypeOf(float64(0))
	propertyTypes[Complex64-1] = reflect.TypeOf(complex64(0))
	propertyTypes[Complex128-1] = reflect.TypeOf(complex128(0))
	propertyTypes[Bytes-1] = reflect.TypeOf([]byte(nil))
	propertyTypes[String-1] = reflect.TypeOf("")
	propertyTypes[Time-1] = reflect.TypeOf(time.Time{})

	propertyTypeOfMap = make(map[reflect.Type]PropertyType, len(propertyTypes))
	for i := PropertyType(1); i < maxPropertyType; i++ {
		propertyTypeOfMap[propertyTypes[i-1]] = i
	}
}
