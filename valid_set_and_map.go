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

	"github.com/donyori/gogo/container"
	"github.com/donyori/gogo/container/mapping"
	"github.com/donyori/gogo/container/set"
	"github.com/donyori/gogo/container/set/mapset"
	"github.com/donyori/gogo/errors"
)

// validSet is a set of items, all of which are valid.
//
// Its method Range accesses items in random order.
// The access order in two calls to Range may be different.
//
// If an invalid item is about to be put into this set,
// the corresponding method panics with the specified error.
//
// To test whether the panic value is of the specified error type,
// convert it to an error with type assertion,
// and then use function errors.As. For example:
//
//	// in a deferred function
//	x := recover()
//	err, ok := x.(error)
//	if ok {
//		var e MyError
//		if errors.As(err, &e) {
//			// x is a MyError
//		}
//	}
type validSet[Item comparable] struct {
	s          set.Set[Item]
	validateFn func(x Item) bool
	errFn      func(x Item) error
}

func _[Item comparable]() {
	var _ set.Set[Item] = (*validSet[Item])(nil)
}

// newValidSet creates a new validSet.
//
// capacity asks to allocate enough space to hold
// the specified number of items.
// If capacity is negative, it is ignored.
//
// validateFn is a function to report whether x is valid.
// If validateFn is nil and x has method "IsValid() bool",
// the validation uses that method.
// If validateFn is nil and x has no such method, newValidSet panics.
//
// errFn is a function returning an error for an invalid item x.
// If errFn is nil, newValidSet uses the following function instead:
//
//	func(x Item) error {
//		return fmt.Errorf("item %v is invalid", x)
//	}
func newValidSet[Item comparable](
	capacity int,
	validateFn func(x Item) bool,
	errFn func(x Item) error,
) *validSet[Item] {
	if validateFn == nil {
		var x Item
		if _, ok := any(x).(interface{ IsValid() bool }); !ok {
			panic(errors.AutoMsg(
				`validateFn is nil and item has no method "IsValid() bool"`))
		}
		validateFn = func(x Item) bool {
			return any(x).(interface{ IsValid() bool }).IsValid()
		}
	}
	if errFn == nil {
		errFn = func(x Item) error {
			return fmt.Errorf("item %v is invalid", x)
		}
	}
	return &validSet[Item]{
		s:          mapset.New[Item](capacity, nil),
		validateFn: validateFn,
		errFn:      errFn,
	}
}

func (vs *validSet[Item]) Len() int {
	return vs.s.Len()
}

// Range accesses the items in the set.
// Each item is accessed once.
// The order of the access is random.
//
// Its parameter handler is a function to deal with the item x in the
// set and report whether to continue to access the next item.
func (vs *validSet[Item]) Range(handler func(x Item) (cont bool)) {
	vs.s.Range(handler)
}

func (vs *validSet[Item]) Filter(filter func(x Item) (keep bool)) {
	vs.s.Filter(filter)
}

func (vs *validSet[Item]) ContainsItem(x Item) bool {
	return vs.s.ContainsItem(x)
}

func (vs *validSet[Item]) ContainsSet(s set.Set[Item]) bool {
	return vs.s.ContainsSet(s)
}

func (vs *validSet[Item]) ContainsAny(c container.Container[Item]) bool {
	return vs.s.ContainsAny(c)
}

func (vs *validSet[Item]) Add(x ...Item) {
	for _, item := range x {
		if !vs.validateFn(item) {
			panic(errors.AutoWrap(vs.errFn(item)))
		}
	}
	vs.s.Add(x...)
}

func (vs *validSet[Item]) Remove(x ...Item) {
	vs.s.Remove(x...)
}

func (vs *validSet[Item]) Union(s set.Set[Item]) {
	if s == nil || s.Len() == 0 {
		return
	}
	vs.validateAllItemsInSet(s)
	vs.s.Union(s)
}

func (vs *validSet[Item]) Intersect(s set.Set[Item]) {
	vs.s.Intersect(s)
}

func (vs *validSet[Item]) Subtract(s set.Set[Item]) {
	vs.s.Subtract(s)
}

func (vs *validSet[Item]) DisjunctiveUnion(s set.Set[Item]) {
	if s == nil || s.Len() == 0 {
		return
	}
	vs.validateAllItemsInSet(s)
	vs.s.DisjunctiveUnion(s)
}

func (vs *validSet[Item]) Clear() {
	vs.s.Clear()
}

// validateAllItemsInSet checks whether all items in s are valid.
//
// If any item is invalid, it panics with the specified error.
//
// The caller should guarantee that s is not nil.
func (vs *validSet[Item]) validateAllItemsInSet(s set.Set[Item]) {
	s.Range(func(x Item) (cont bool) {
		if !vs.validateFn(x) {
			panic(errors.AutoWrapSkip(vs.errFn(x), 2))
		}
		return true
	})
}

// validMap is a map in which all keys and values are valid.
//
// Its method Range accesses key-value pairs in random order.
// The access order in two calls to Range may be different.
//
// If an invalid key or value is about to be put into this map,
// the corresponding method panics with the specified error.
//
// To test whether the panic value is of the specified error type,
// convert it to an error with type assertion,
// and then use function errors.As. For example:
//
//	// in a deferred function
//	x := recover()
//	err, ok := x.(error)
//	if ok {
//		var e MyError
//		if errors.As(err, &e) {
//			// x is a MyError
//		}
//	}
type validMap[Key comparable, Value any] struct {
	m               mapping.GoMap[Key, Value]
	keyValidateFn   func(key Key) bool
	keyErrFn        func(key Key) error
	valueValidateFn func(value Value) bool
	valueErrFn      func(value Value) error
}

func _[Key comparable, Value any]() {
	var _ mapping.Map[Key, Value] = (*validMap[Key, Value])(nil)
}

// newValidMap creates a new validMap.
//
// capacity asks to allocate enough space to hold
// the specified number of key-value pairs.
// If capacity is negative, it is ignored.
//
// keyValidateFn and valueValidateFn are functions to report
// whether the key and value are valid, respectively.
// If keyValidateFn is nil and the key has method "IsValid() bool",
// the validation uses that method.
// If keyValidateFn is nil and the key has no such method, newValidMap panics.
// Similarly, if valueValidateFn is nil and the value has method
// "IsValid() bool", the validation uses that method.
// If valueValidateFn is nil and the value has no such method,
// newValidMap panics.
//
// keyErrFn and valueErrFn are functions returning errors
// for an invalid key and value, respectively.
// If keyErrFn is nil, newValidMap uses the following function instead:
//
//	func(key Key) error {
//		return fmt.Errorf("key %v is invalid", key)
//	}
//
// Similarly, if valueErrFn is nil,
// newValidMap uses the following function instead:
//
//	func(value Value) error {
//		return fmt.Errorf("value %v is invalid", value)
//	}
func newValidMap[Key comparable, Value any](
	capacity int,
	keyValidateFn func(key Key) bool,
	keyErrFn func(key Key) error,
	valueValidateFn func(value Value) bool,
	valueErrFn func(value Value) error,
) *validMap[Key, Value] {
	var m mapping.GoMap[Key, Value]
	if capacity > 0 {
		m = make(mapping.GoMap[Key, Value], capacity)
	}
	if keyValidateFn == nil {
		var k Key
		if _, ok := any(k).(interface{ IsValid() bool }); !ok {
			panic(errors.AutoMsg(
				`keyValidateFn is nil and key has no method "IsValid() bool"`))
		}
		keyValidateFn = func(key Key) bool {
			return any(key).(interface{ IsValid() bool }).IsValid()
		}
	}
	if keyErrFn == nil {
		keyErrFn = func(key Key) error {
			return fmt.Errorf("key %v is invalid", key)
		}
	}
	if valueValidateFn == nil {
		var v Value
		if _, ok := any(v).(interface{ IsValid() bool }); !ok {
			panic(errors.AutoMsg(
				`valueValidateFn is nil and value has no method "IsValid() bool"`))
		}
		valueValidateFn = func(value Value) bool {
			return any(v).(interface{ IsValid() bool }).IsValid()
		}
	}
	if valueErrFn == nil {
		valueErrFn = func(value Value) error {
			return fmt.Errorf("value %v is invalid", value)
		}
	}
	return &validMap[Key, Value]{
		m:               m,
		keyValidateFn:   keyValidateFn,
		keyErrFn:        keyErrFn,
		valueValidateFn: valueValidateFn,
		valueErrFn:      valueErrFn,
	}
}

func (vm *validMap[Key, Value]) Len() int {
	return vm.m.Len()
}

// Range accesses the key-value pairs in the map.
// Each key-value pair is accessed once.
// The order of the access is random.
//
// Its parameter handler is a function to deal with the key-value pair x
// in the map and report whether to continue to access the next key-value pair.
func (vm *validMap[Key, Value]) Range(
	handler func(x mapping.Entry[Key, Value]) (cont bool)) {
	vm.m.Range(handler)
}

func (vm *validMap[Key, Value]) Filter(
	filter func(x mapping.Entry[Key, Value]) (keep bool)) {
	vm.m.Filter(filter)
}

func (vm *validMap[Key, Value]) Get(key Key) (value Value, present bool) {
	return vm.m.Get(key)
}

func (vm *validMap[Key, Value]) Set(key Key, value Value) {
	vm.validateKeyAndValue(key, value)
	vm.m.Set(key, value)
}

func (vm *validMap[Key, Value]) GetAndSet(key Key, value Value) (
	previous Value, present bool) {
	vm.validateKeyAndValue(key, value)
	return vm.m.GetAndSet(key, value)
}

func (vm *validMap[Key, Value]) SetMap(m mapping.Map[Key, Value]) {
	if m == nil || m.Len() == 0 {
		return
	}
	vm.validateAllKeysAndValuesInMap(m)
	vm.m.SetMap(m)
}

func (vm *validMap[Key, Value]) GetAndSetMap(m mapping.Map[Key, Value]) (
	previous mapping.Map[Key, Value]) {
	if m == nil || m.Len() == 0 {
		return
	}
	vm.validateAllKeysAndValuesInMap(m)
	return vm.m.GetAndSetMap(m)
}

func (vm *validMap[Key, Value]) Remove(key ...Key) {
	vm.m.Remove(key...)
}

func (vm *validMap[Key, Value]) GetAndRemove(key Key) (
	previous Value, present bool) {
	return vm.m.GetAndRemove(key)
}

func (vm *validMap[Key, Value]) Clear() {
	vm.m.Clear()
}

// validateKeyAndValue checks whether key and value are valid.
//
// If not, it panics with the specified error.
func (vm *validMap[Key, Value]) validateKeyAndValue(key Key, value Value) {
	if !vm.keyValidateFn(key) {
		panic(errors.AutoWrapSkip(vm.keyErrFn(key), 1))
	} else if !vm.valueValidateFn(value) {
		panic(errors.AutoWrapSkip(vm.valueErrFn(value), 1))
	}
}

// validateAllKeysAndValuesInMap checks whether
// all key-value pairs in m are valid.
//
// If not, it panics with the specified error.
//
// The caller should guarantee that m is not nil.
func (vm *validMap[Key, Value]) validateAllKeysAndValuesInMap(
	m mapping.Map[Key, Value]) {
	m.Range(func(x mapping.Entry[Key, Value]) (cont bool) {
		if !vm.keyValidateFn(x.Key) {
			panic(errors.AutoWrapSkip(vm.keyErrFn(x.Key), 2))
		} else if !vm.valueValidateFn(x.Value) {
			panic(errors.AutoWrapSkip(vm.valueErrFn(x.Value), 2))
		}
		return true
	})
}
