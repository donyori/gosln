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
	"strings"

	"github.com/donyori/gogo/container"
	"github.com/donyori/gogo/container/set"
	"github.com/donyori/gogo/errors"
)

// encode64Table is a character table for encoding
// the serial number (int64) to a valid suffix of ID.
const encode64Table = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

// IsValidTypeString reports whether t is a valid type value.
//
// A valid type consists of alphanumeric characters and underscores ('_'),
// begins with an uppercase letter, does not begin with "SLN",
// and is up to 65535 bytes long.
func IsValidTypeString(t string) bool {
	if len(t) < 1 || len(t) > 65535 ||
		t[0] < 'A' || t[0] > 'Z' ||
		len(t) >= 3 && t[:3] == "SLN" {
		return false
	}
	for i := 1; i < len(t); i++ {
		if t[i] < '0' ||
			t[i] > 'z' ||
			t[i] > '9' && t[i] < 'A' ||
			t[i] > 'Z' && t[i] != '_' && t[i] < 'a' {
			return false
		}
	}
	return true
}

// Type is the type of the semantic node and link.
//
// A valid type consists of alphanumeric characters and underscores ('_'),
// begins with an uppercase letter, does not begin with "SLN",
// and is up to 65535 bytes long.
type Type struct {
	t string
}

// NewType returns a Type whose value is t.
//
// A valid type consists of alphanumeric characters and underscores ('_'),
// begins with an uppercase letter, does not begin with "SLN",
// and is up to 65535 bytes long.
// If t is invalid, NewType reports a *InvalidTypeError.
// (To test whether err is *InvalidTypeError, use function errors.As.)
func NewType(t string) (typ Type, err error) {
	if IsValidTypeString(t) {
		typ.t = t
	} else {
		err = errors.AutoWrap(NewInvalidTypeError(t))
	}
	return
}

// MustNewType is like NewType, but panic if the type is invalid.
func MustNewType(t string) Type {
	typ, err := NewType(t)
	if err != nil {
		panic(errors.AutoWrap(err))
	}
	return typ
}

// String returns the type value.
//
// If t is invalid, String returns an empty string.
func (t Type) String() string {
	return t.t
}

// IsValid reports whether t is valid.
func (t Type) IsValid() bool {
	// Its constructor should guarantee that t is valid if it is not zero.
	return t.t != ""
}

// ID is the unique identifier of the semantic node and link.
//
// A valid ID is the concatenation of its corresponding type,
// a number sign ('#'), and a unique suffix.
//
// ID should be assigned by the Semantic Link Network, not by the client.
type ID struct {
	t string // The corresponding type.
	s string // The suffix, unique across the IDs for the same type.
}

// NewID returns an ID corresponding to the type t,
// with the specified date and the serial number i.
//
// If t is invalid (such as zero-value),
// NewID returns a zero-value ID.
//
// If i is negative, NewID panics.
func NewID(t Type, date Date, i int64) ID {
	if i < 0 {
		panic(errors.AutoMsg(fmt.Sprintf("the number i (%d) is negative", i)))
	}
	if !t.IsValid() {
		return ID{}
	}
	var b strings.Builder
	b.Grow(19)
	b.WriteString(date.String())
	b.WriteByte('-')
	for {
		b.WriteByte(encode64Table[i&077])
		i >>= 6
		if i == 0 {
			return ID{
				t: t.String(),
				s: b.String(),
			}
		}
		i--
	}
}

// String formats id into a string in the form of
//
//	<Type> "#" <UniqueSuffix>
//
// where <Type> is the type corresponding to id,
// and <UniqueSuffix> is a suffix that is unique
// across the IDs for the same type.
//
// If id is invalid, String returns an empty string.
func (id ID) String() string {
	if id.t == "" {
		return ""
	}
	return id.t + "#" + id.s
}

// IsValid reports whether id is valid.
func (id ID) IsValid() bool {
	// Its constructor should guarantee that id is valid if it is not zero.
	return id.t != ""
}

// Type returns the type corresponding to id.
func (id ID) Type() Type {
	if id.t == "" {
		return Type{}
	}
	return MustNewType(id.t)
}

// TypeSet is a set of node or link types, all of which are valid Type.
//
// If an invalid Type is about to be put into this set,
// the corresponding method panics with a *InvalidTypeError.
//
// To test whether the panic value is a *InvalidTypeError,
// convert it to an error with type assertion,
// and then use function errors.As. For example:
//
//	// in a deferred function
//	x := recover()
//	err, ok := x.(error)
//	if ok {
//		var e *gosln.InvalidTypeError
//		if errors.As(err, &e) {
//			// x is a *InvalidTypeError
//		}
//	}
type TypeSet interface {
	set.Set[Type]
}

// NewTypeSet creates a new TypeSet.
//
// The method Range of the set accesses types in random order.
// The access order in two calls to Range may be different.
//
// capacity asks to allocate enough space to hold
// the specified number of types.
// If capacity is negative, it is ignored.
func NewTypeSet(capacity int) TypeSet {
	return newValidSet(
		capacity,
		func(x Type) bool {
			return x.IsValid()
		},
		func(x Type) error {
			return NewInvalidTypeError(x.String())
		},
	)
}

// IDSet is a set of IDs, where the IDs are valid.
//
// If an invalid ID is about to be put into this set,
// the corresponding method panics with a *InvalidIDError.
//
// To test whether the panic value is a *InvalidIDError,
// convert it to an error with type assertion,
// and then use function errors.As. For example:
//
//	// in a deferred function
//	x := recover()
//	err, ok := x.(error)
//	if ok {
//		var e *gosln.InvalidIDError
//		if errors.As(err, &e) {
//			// x is a *InvalidIDError
//		}
//	}
type IDSet interface {
	set.Set[ID]

	// LenType returns the number of IDs
	// corresponding to the type t in the set.
	LenType(t Type) int

	// NumType returns the number of types
	// corresponding to the IDs in the set.
	NumType() int

	// RangeType accesses the IDs corresponding to the type t in the set.
	// Each ID is accessed once. The order of the access is random.
	//
	// Its parameter handler is a function to deal with an ID in the set
	// and report whether to continue to access the next ID.
	RangeType(t Type, handler func(id ID) (cont bool))

	// ContainsType reports whether there is an ID
	// corresponding to the type t in the set.
	ContainsType(t Type) bool
}

// idSetImpl is an implementation of interface IDSet.
type idSetImpl struct {
	m map[string]map[string]struct{}
}

// NewIDSet creates a new IDSet.
//
// The method Range of the set accesses IDs in random order.
// The access order in two calls to Range may be different.
func NewIDSet() IDSet {
	return &idSetImpl{m: make(map[string]map[string]struct{})}
}

func (ids *idSetImpl) Len() int {
	var n int
	for _, sub := range ids.m {
		n += len(sub)
	}
	return n
}

// Range accesses the IDs in the set.
// Each ID is accessed once.
// The order of the access is random.
//
// Its parameter handler is a function to deal with an ID in the set
// and report whether to continue to access the next ID.
func (ids *idSetImpl) Range(handler func(x ID) (cont bool)) {
	for t, sub := range ids.m {
		for suffix := range sub {
			if !handler(ID{t: t, s: suffix}) {
				return
			}
		}
	}
}

func (ids *idSetImpl) Filter(filter func(x ID) (keep bool)) {
	for t, sub := range ids.m {
		for suffix := range sub {
			if !filter(ID{t: t, s: suffix}) {
				delete(sub, suffix)
				if len(sub) == 0 {
					delete(ids.m, t)
				}
			}
		}
	}
}

func (ids *idSetImpl) ContainsItem(x ID) bool {
	sub := ids.m[x.t]
	if sub == nil {
		return false
	}
	_, ok := sub[x.s]
	return ok
}

func (ids *idSetImpl) ContainsSet(s set.Set[ID]) bool {
	if s == nil {
		return true
	}
	n := s.Len()
	if n == 0 {
		return true
	} else if n > len(ids.m) {
		return false
	}
	var ok bool
	s.Range(func(x ID) (cont bool) {
		sub := ids.m[x.t]
		if sub != nil {
			_, ok = sub[x.s]
		} else {
			ok = false
		}
		return ok
	})
	return ok
}

func (ids *idSetImpl) ContainsAny(c container.Container[ID]) bool {
	if c == nil || c.Len() == 0 {
		return false
	}
	var ok bool
	c.Range(func(x ID) (cont bool) {
		sub := ids.m[x.t]
		if sub != nil {
			_, ok = sub[x.s]
		} else {
			ok = false
		}
		return !ok
	})
	return ok
}

func (ids *idSetImpl) Add(id ...ID) {
	for _, x := range id {
		if !x.IsValid() {
			panic(errors.AutoWrap(NewInvalidIDError(x)))
		}
	}
	for _, x := range id {
		sub := ids.m[x.t]
		if sub == nil {
			sub = make(map[string]struct{})
			ids.m[x.t] = sub
		}
		sub[x.s] = struct{}{}
	}
}

func (ids *idSetImpl) Remove(id ...ID) {
	for _, x := range id {
		sub := ids.m[x.t]
		if sub != nil {
			delete(sub, x.s)
			if len(sub) == 0 {
				delete(ids.m, x.t)
			}
		}
	}
}

func (ids *idSetImpl) Union(s set.Set[ID]) {
	if s == nil || s.Len() == 0 {
		return
	}
	validateAllIDsInSet(s)
	s.Range(func(x ID) (cont bool) {
		sub := ids.m[x.t]
		if sub == nil {
			sub = make(map[string]struct{})
			ids.m[x.t] = sub
		}
		sub[x.s] = struct{}{}
		return true
	})
}

func (ids *idSetImpl) Intersect(s set.Set[ID]) {
	if s == nil || s.Len() == 0 {
		ids.m = make(map[string]map[string]struct{})
		return
	}
	for t, sub := range ids.m {
		for suffix := range sub {
			if !s.ContainsItem(ID{t: t, s: suffix}) {
				delete(sub, suffix)
				if len(sub) == 0 {
					delete(ids.m, t)
				}
			}
		}
	}
}

func (ids *idSetImpl) Subtract(s set.Set[ID]) {
	if s == nil || s.Len() == 0 {
		return
	}
	s.Range(func(x ID) (cont bool) {
		sub := ids.m[x.t]
		if sub != nil {
			delete(sub, x.s)
			if len(sub) == 0 {
				delete(ids.m, x.t)
			}
		}
		return true
	})
}

func (ids *idSetImpl) DisjunctiveUnion(s set.Set[ID]) {
	if s == nil || s.Len() == 0 {
		return
	}
	validateAllIDsInSet(s)
	s.Range(func(x ID) (cont bool) {
		sub := ids.m[x.t]
		if sub == nil {
			sub = make(map[string]struct{})
			ids.m[x.t] = sub
		}
		if _, ok := sub[x.s]; ok {
			delete(sub, x.s)
			if len(sub) == 0 {
				delete(ids.m, x.t)
			}
		} else {
			sub[x.s] = struct{}{}
		}
		return true
	})
}

func (ids *idSetImpl) Clear() {
	ids.m = make(map[string]map[string]struct{})
}

func (ids *idSetImpl) LenType(t Type) int {
	return len(ids.m[t.t])
}

func (ids *idSetImpl) NumType() int {
	return len(ids.m)
}

func (ids *idSetImpl) RangeType(t Type, handler func(id ID) (cont bool)) {
	for suffix := range ids.m[t.t] {
		if !handler(ID{t: t.t, s: suffix}) {
			return
		}
	}
}

func (ids *idSetImpl) ContainsType(t Type) bool {
	return len(ids.m[t.t]) > 0
}

// validateAllIDsInSet checks whether all IDs in s are valid.
//
// If any ID is invalid, it panics with a *InvalidIDError.
func validateAllIDsInSet(s set.Set[ID]) {
	if s == nil {
		return
	}
	s.Range(func(x ID) (cont bool) {
		if !x.IsValid() {
			panic(errors.AutoWrapSkip(NewInvalidIDError(x), 2))
		}
		return true
	})
}
