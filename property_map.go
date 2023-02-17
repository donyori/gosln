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
	"regexp"
	"time"

	"github.com/donyori/gogo/constraints"
	"github.com/donyori/gogo/errors"
)

// PropertyValue is a constraint for property values
// of semantic nodes and links.
//
// It matches the following types:
//   - Built-in Boolean: bool.
//   - Built-in signed integers: int, int8, int16, int32 (rune), int64.
//   - Built-in unsigned integers: uint, uint8 (byte), uint16, uint32, uint64, uintptr.
//   - Built-in floating numbers: float32, float64.
//   - Built-in complex numbers: complex64, complex128.
//   - Byte strings: []byte, string.
//   - Time: time.Time.
type PropertyValue interface {
	bool |
		constraints.PredeclaredNumeric |
		constraints.PredeclaredByteString |
		time.Time
}

// PropertyMap is a property name-value map,
// where the name is a string consisting of alphanumeric characters
// and underscores ('_'), beginning with a lowercase letter,
// and up to 65535 bytes long,
// and the value type should conform to PropertyValue.
//
// A zero-value PropertyMap is ready to use.
type PropertyMap struct {
	m map[string]any
}

// Len returns the number of properties in the map.
func (pm *PropertyMap) Len() int {
	if pm == nil {
		return 0
	}
	return len(pm.m)
}

// Get obtains the property with the specified name and
// returns its type and value.
//
// If the property does not exist, it will return (0, nil).
// To test whether the property exists,
// it is sufficient to test whether the property type t is 0.
//
// To get a property of a specified type,
// use the generic function GetProperty.
func (pm *PropertyMap) Get(name string) (t PropertyType, value any) {
	if pm == nil || len(pm.m) == 0 {
		return
	}
	value, ok := pm.m[name]
	if ok {
		t = PropertyTypeOf(value)
	}
	return
}

// Range accesses the properties in the map.
// Each property will be accessed once.
// The access order may be random and may be different at each call.
//
// Its parameter handler is a function to deal with the property
// with the specified name, type, and value in the map and
// report whether to continue to access the next property.
func (pm *PropertyMap) Range(handler func(name string, t PropertyType, value any) (cont bool)) {
	if pm != nil {
		for name, value := range pm.m {
			if !handler(name, PropertyTypeOf(value), value) {
				return
			}
		}
	}
}

// Remove removes the property with the specified name in the map and
// returns its type and value.
//
// If the property does not exist, it will do nothing and return (0, nil).
// To test whether the property exists,
// it is sufficient to test whether the property type t is 0.
func (pm *PropertyMap) Remove(name string) (t PropertyType, value any) {
	if pm == nil || len(pm.m) == 0 {
		return
	}
	value, ok := pm.m[name]
	if ok {
		t = PropertyTypeOf(value)
		delete(pm.m, name)
	}
	return
}

// Clear removes all properties in the map.
func (pm *PropertyMap) Clear() {
	if pm != nil {
		pm.m = nil
	}
}

// GetProperty obtains the property with the specified name from pm.
//
// If the property does not exist, it will report a *PropertyNotExistError.
// If the type of the property is not V and not convertible to V,
// it will report a *PropertyTypeError.
// (To test the type of err, use function errors.As.)
func GetProperty[V PropertyValue](pm *PropertyMap, name string) (value V, err error) {
	if pm == nil || len(pm.m) == 0 {
		err = errors.AutoWrap(NewPropertyNotExistError(name))
		return
	}
	prop, ok := pm.m[name]
	if !ok {
		err = errors.AutoWrap(NewPropertyNotExistError(name))
		return
	}
	propV := reflect.ValueOf(prop)
	// Call ValueOf on the pointer of value so that
	// the value can be settable for basic types.
	v := reflect.ValueOf(&value).Elem()
	propType, vType := propV.Type(), v.Type()
	switch {
	case propType.AssignableTo(vType):
		v.Set(propV)
	case propType.ConvertibleTo(vType):
		v.Set(propV.Convert(vType))
	default:
		err = errors.AutoWrap(
			NewPropertyTypeError(name, prop, vType.String()),
		)
	}
	return
}

// propertyNamePattern is the regular expression pattern for property names.
var propertyNamePattern = regexp.MustCompile("^[a-z][0-9A-Z_a-z]{0,65534}$")

// SetProperty sets a property with the specified name and value to pm.
//
// A valid property name consists of alphanumeric characters and
// underscores ('_'), begins with a lowercase letter,
// and is up to 65535 bytes long.
//
// If pm is nil, it will report an error.
// If name is invalid, it will report a *InvalidPropertyNameError.
// (To test whether the error is *InvalidPropertyNameError,
// use function errors.As.)
func SetProperty[V PropertyValue](pm *PropertyMap, name string, value V) error {
	switch {
	case pm == nil:
		return errors.AutoNew("property map is nil")
	case !propertyNamePattern.MatchString(name):
		return errors.AutoWrap(NewInvalidPropertyNameError(name))
	case pm.m == nil:
		pm.m = make(map[string]any)
	}
	pm.m[name] = value
	return nil
}
