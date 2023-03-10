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

	"github.com/donyori/gogo/container/mapping"
)

// PropType represents the type of property.
type PropType int8

const (
	PTBool       PropType = 1 + iota // bool
	PTInt                            // int
	PTInt8                           // int8
	PTInt16                          // int16
	PTInt32                          // int32
	PTInt64                          // int64
	PTUint                           // uint
	PTUint8                          // uint8
	PTUint16                         // uint16
	PTUint32                         // uint32
	PTUint64                         // uint64
	PTUintptr                        // uintptr
	PTFloat32                        // float32
	PTFloat64                        // float64
	PTComplex64                      // complex64
	PTComplex128                     // complex128
	PTBytes                          // []byte
	PTString                         // string
	PTTime                           // time.Time
	PTDate                           // gosln.Date
	maxPropType                      // PropType(21)
)

// Before running the following command, please make sure the numeric value
// in the line comment of maxPropType is correct.
//
//go:generate stringer -type=PropType -output=prop_type_string.go -linecomment

var (
	// propTypes is an array of reflect.Type corresponding to PropType.
	propTypes [maxPropType - 1]reflect.Type
	// propTypeOfMap is a map from reflect.Type to PropType,
	//used by PropTypeOf.
	propTypeOfMap map[reflect.Type]PropType
)

func init() {
	propTypes[PTBool-1] = reflect.TypeOf(false)
	propTypes[PTInt-1] = reflect.TypeOf(0)
	propTypes[PTInt8-1] = reflect.TypeOf(int8(0))
	propTypes[PTInt16-1] = reflect.TypeOf(int16(0))
	propTypes[PTInt32-1] = reflect.TypeOf(int32(0))
	propTypes[PTInt64-1] = reflect.TypeOf(int64(0))
	propTypes[PTUint-1] = reflect.TypeOf(uint(0))
	propTypes[PTUint8-1] = reflect.TypeOf(uint8(0))
	propTypes[PTUint16-1] = reflect.TypeOf(uint16(0))
	propTypes[PTUint32-1] = reflect.TypeOf(uint32(0))
	propTypes[PTUint64-1] = reflect.TypeOf(uint64(0))
	propTypes[PTUintptr-1] = reflect.TypeOf(uintptr(0))
	propTypes[PTFloat32-1] = reflect.TypeOf(float32(0))
	propTypes[PTFloat64-1] = reflect.TypeOf(float64(0))
	propTypes[PTComplex64-1] = reflect.TypeOf(complex64(0))
	propTypes[PTComplex128-1] = reflect.TypeOf(complex128(0))
	propTypes[PTBytes-1] = reflect.TypeOf([]byte(nil))
	propTypes[PTString-1] = reflect.TypeOf("")
	propTypes[PTTime-1] = reflect.TypeOf(time.Time{})
	propTypes[PTDate-1] = reflect.TypeOf(Date{})

	propTypeOfMap = make(map[reflect.Type]PropType, len(propTypes))
	for i := PropType(1); i < maxPropType; i++ {
		propTypeOfMap[propTypes[i-1]] = i
	}
}

// PropTypeOf returns the property type of the value v.
//
// It returns 0 if v does not conform to PropValue.
func PropTypeOf(v any) PropType {
	return propTypeOfMap[reflect.TypeOf(v)]
}

// IsValid reports whether the property type is known.
func (i PropType) IsValid() bool {
	return i > 0 && i < maxPropType
}

// GoType returns the reflect.Type corresponding to the property type.
//
// It returns nil if the property type is invalid.
func (i PropType) GoType() reflect.Type {
	if i > 0 && i < maxPropType {
		return propTypes[i-1]
	}
	return nil
}

// IsConvertibleTo reports whether the property type i can convert to type t.
func (i PropType) IsConvertibleTo(t PropType) bool {
	if i <= 0 || i >= maxPropType || t <= 0 || t >= maxPropType {
		return false
	}
	return propTypes[i-1].ConvertibleTo(propTypes[t-1])
}

// IsSignedInteger reports whether the property type is a signed integer,
// including int, int8, int16, int32 (rune), and int64.
func (i PropType) IsSignedInteger() bool {
	switch i {
	case PTInt, PTInt8, PTInt16, PTInt32, PTInt64:
		return true
	}
	return false
}

// IsUnsignedInteger reports whether the property type is an unsigned integer,
// including uint, uint8 (byte), uint16, uint32, uint64, and uintptr.
func (i PropType) IsUnsignedInteger() bool {
	switch i {
	case PTUint, PTUint8, PTUint16, PTUint32, PTUint64, PTUintptr:
		return true
	}
	return false
}

// IsInteger reports whether the property type is an integer,
// including int, int8, int16, int32 (rune), int64,
// uint, uint8 (byte), uint16, uint32, uint64, and uintptr.
func (i PropType) IsInteger() bool {
	switch i {
	case PTInt, PTInt8, PTInt16, PTInt32, PTInt64,
		PTUint, PTUint8, PTUint16, PTUint32, PTUint64, PTUintptr:
		return true
	}
	return false
}

// IsFloat reports whether the property type is a floating number,
// including float32 and float64.
func (i PropType) IsFloat() bool {
	switch i {
	case PTFloat32, PTFloat64:
		return true
	}
	return false
}

// IsRealNumber reports whether the property type is a real number,
// including int, int8, int16, int32 (rune), int64,
// uint, uint8 (byte), uint16, uint32, uint64, uintptr,
// float32, and float64.
func (i PropType) IsRealNumber() bool {
	switch i {
	case PTInt, PTInt8, PTInt16, PTInt32, PTInt64,
		PTUint, PTUint8, PTUint16, PTUint32, PTUint64, PTUintptr,
		PTFloat32, PTFloat64:
		return true
	}
	return false
}

// IsComplex reports whether the property type is a complex number,
// including complex64 and complex128.
func (i PropType) IsComplex() bool {
	switch i {
	case PTComplex64, PTComplex128:
		return true
	}
	return false
}

// IsNumeric reports whether the property type is a number,
// including int, int8, int16, int32 (rune), int64,
// uint, uint8 (byte), uint16, uint32, uint64, uintptr,
// float32, float64, complex64, and complex128.
func (i PropType) IsNumeric() bool {
	switch i {
	case PTInt, PTInt8, PTInt16, PTInt32, PTInt64,
		PTUint, PTUint8, PTUint16, PTUint32, PTUint64, PTUintptr,
		PTFloat32, PTFloat64,
		PTComplex64, PTComplex128:
		return true
	}
	return false
}

// IsByteString reports whether the property type is a byte string,
// including []byte and string.
func (i PropType) IsByteString() bool {
	switch i {
	case PTBytes, PTString:
		return true
	}
	return false
}

// PropTypeMap is a property name-type map,
// where the names are valid PropName
// and the types are valid PropType.
//
// If an invalid PropName is about to be put into this map,
// the corresponding method panics with a *InvalidPropNameError.
//
// If an invalid PropType is about to be put into this map,
// the corresponding method panics with a *InvalidPropTypeError.
//
// To test whether the panic value is a *InvalidPropNameError or
// *InvalidPropTypeError, convert it to an error with type assertion,
// and then use function errors.As. For example:
//
//	// in a deferred function
//	x := recover()
//	err, ok := x.(error)
//	if ok {
//		var e *gosln.InvalidPropNameError
//		if errors.As(err, &e) {
//			// x is a *InvalidPropNameError
//		}
//	}
type PropTypeMap interface {
	mapping.Map[PropName, PropType]
}

// NewPropTypeMap creates a new PropTypeMap.
//
// The method Range of the map accesses
// property name-type pairs in random order.
// The access order in two calls to Range may be different.
//
// capacity asks to allocate enough space to hold
// the specified number of property name-type pairs.
// If capacity is negative, it is ignored.
func NewPropTypeMap(capacity int) PropTypeMap {
	return newValidMap(
		capacity,
		func(key PropName) bool {
			return key.IsValid()
		},
		func(key PropName) error {
			return NewInvalidPropNameError(key.String())
		},
		func(value PropType) bool {
			return value.IsValid()
		},
		func(value PropType) error {
			return NewInvalidPropTypeError(value)
		},
	)
}
