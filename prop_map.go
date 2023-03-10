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

	"github.com/donyori/gogo/constraints"
	"github.com/donyori/gogo/container/mapping"
	"github.com/donyori/gogo/errors"
)

// PropValue is a constraint for property values
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
//   - Date: gosln.Date.
type PropValue interface {
	bool |
		constraints.PredeclaredNumeric |
		constraints.PredeclaredByteString |
		time.Time | Date
}

// PropMap is a property name-value map,
// where the names are valid PropName
// and the value types conform to PropValue.
//
// If an invalid PropName is about to be put into this map,
// the corresponding method panics with a *InvalidPropNameError.
//
// If a property value that does not conform to PropValue
// is about to be put into this map,
// the corresponding method panics with a *InvalidPropValueError.
//
// To test whether the panic value is a *InvalidPropNameError or
// *InvalidPropValueError, convert it to an error with type assertion,
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
type PropMap interface {
	mapping.Map[PropName, any]
}

// NewPropMap creates a new PropMap.
//
// The method Range of the map accesses properties in random order.
// The access order in two calls to Range may be different.
//
// capacity asks to allocate enough space to hold
// the specified number of properties.
// If capacity is negative, it is ignored.
func NewPropMap(capacity int) PropMap {
	return newValidMap(
		capacity,
		func(key PropName) bool {
			return key.IsValid()
		},
		func(key PropName) error {
			return NewInvalidPropNameError(key.String())
		},
		func(value any) bool {
			return PropTypeOf(value).IsValid()
		},
		func(value any) error {
			return NewInvalidPropValueError(value)
		},
	)
}

// mutExclPropMap is an implementation of interface PropMap.
//
// It can associate with one or more collections
// that have the method Remove(...PropName).
// When a property is put into this map,
// mutExclPropMap removes the property name from these collections.
//
// The client must call its method init to initialize
// the mutExclPropMap before use.
type mutExclPropMap struct {
	m PropMap
	r []interface{ Remove(...PropName) }
}

// init initializes the mutExclPropMap
// with the specified capacity and collections.
//
// capacity asks to allocate enough space to hold
// the specified number of properties.
// If capacity is negative, it is ignored.
//
// collection is a list of collections associated with this map.
// When a property is put into this map,
// mutExclPropMap removes the property name from these collections.
func (mepm *mutExclPropMap) init(capacity int,
	collection ...interface{ Remove(...PropName) }) {
	mepm.m = NewPropMap(capacity)
	if len(collection) > 0 {
		mepm.r = make([]interface{ Remove(...PropName) }, len(collection))
		copy(mepm.r, collection)
	}
}

func (mepm *mutExclPropMap) Len() int {
	mepm.checkInit()
	return mepm.m.Len()
}

// Range accesses the properties in the map.
// Each property is accessed once.
// The access order may be random and may be different at each call.
//
// Its parameter handler is a function to deal with the property
// with the specified name and value in the map and
// report whether to continue to access the next property.
func (mepm *mutExclPropMap) Range(
	handler func(x mapping.Entry[PropName, any]) (cont bool)) {
	mepm.checkInit()
	mepm.m.Range(handler)
}

func (mepm *mutExclPropMap) Filter(
	filter func(x mapping.Entry[PropName, any]) (keep bool)) {
	mepm.checkInit()
	mepm.m.Filter(filter)
}

func (mepm *mutExclPropMap) Get(key PropName) (value any, present bool) {
	mepm.checkInit()
	return mepm.m.Get(key)
}

func (mepm *mutExclPropMap) Set(key PropName, value any) {
	mepm.checkInit()
	mepm.m.Set(key, value)
	mepm.removeFromOthers(key)
}

func (mepm *mutExclPropMap) GetAndSet(key PropName, value any) (
	previous any, present bool) {
	mepm.checkInit()
	previous, present = mepm.m.GetAndSet(key, value)
	mepm.removeFromOthers(key)
	return
}

func (mepm *mutExclPropMap) SetMap(m mapping.Map[PropName, any]) {
	mepm.checkInit()
	if m == nil || m.Len() == 0 {
		return
	}
	mepm.m.SetMap(m)
	m.Range(func(x mapping.Entry[PropName, any]) (cont bool) {
		mepm.removeFromOthers(x.Key)
		return true
	})
}

func (mepm *mutExclPropMap) GetAndSetMap(m mapping.Map[PropName, any]) (
	previous mapping.Map[PropName, any]) {
	mepm.checkInit()
	if m == nil || m.Len() == 0 {
		return
	}
	previous = mepm.m.GetAndSetMap(m)
	m.Range(func(x mapping.Entry[PropName, any]) (cont bool) {
		mepm.removeFromOthers(x.Key)
		return true
	})
	return
}

func (mepm *mutExclPropMap) Remove(key ...PropName) {
	mepm.checkInit()
	mepm.m.Remove(key...)
}

func (mepm *mutExclPropMap) GetAndRemove(key PropName) (
	previous any, present bool) {
	mepm.checkInit()
	return mepm.m.GetAndRemove(key)
}

func (mepm *mutExclPropMap) Clear() {
	mepm.checkInit()
	mepm.m.Clear()
}

// checkInit checks whether mepm is initialized.
// If not, it panics.
func (mepm *mutExclPropMap) checkInit() {
	if mepm.m == nil {
		panic(errors.AutoMsgCustom("not initialized before use", -1, 1))
	}
}

// removeFromOthers removes name from collections in mepm.r.
func (mepm *mutExclPropMap) removeFromOthers(name ...PropName) {
	if len(name) > 0 {
		for _, r := range mepm.r {
			r.Remove(name...)
		}
	}
}

// PropMapGet obtains the property with the specified name from pm.
//
// If the property does not exist, it reports a *PropNotExistError.
// If the type of the property is not V and not convertible to V,
// it reports a *PropTypeError.
// (To test the type of err, use function errors.As.)
func PropMapGet[V PropValue](pm PropMap, name PropName) (value V, err error) {
	if pm == nil {
		err = errors.AutoWrap(NewPropNotExistError(name))
		return
	}
	prop, present := pm.Get(name)
	if !present {
		err = errors.AutoWrap(NewPropNotExistError(name))
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
		err = errors.AutoWrap(NewPropTypeError(name, prop, vType))
	}
	return
}

// PropMapSet sets a property with the specified name and value to pm.
//
// If pm is nil, it reports an error.
// If name is invalid, it reports a *InvalidPropNameError.
// (To test whether the error is *InvalidPropNameError,
// use function errors.As.)
func PropMapSet[V PropValue](pm PropMap, name PropName, value V) error {
	if pm == nil {
		return errors.AutoNew("property map is nil")
	} else if !name.IsValid() {
		return errors.AutoWrap(NewInvalidPropNameError(name.String()))
	}
	pm.Set(name, value)
	return nil
}
