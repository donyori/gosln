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
	"github.com/donyori/gogo/container"
	"github.com/donyori/gogo/container/set"
	"github.com/donyori/gogo/errors"
)

// IsValidPropNameString reports whether name is a valid property name.
//
// A valid property name consists of alphanumeric characters and
// underscores ('_'), begins with a lowercase letter,
// does not begin with "sln", and is up to 65535 bytes long.
func IsValidPropNameString(name string) bool {
	if len(name) < 1 || len(name) > 65535 ||
		name[0] < 'a' || name[0] > 'z' ||
		len(name) >= 3 && name[:3] == "sln" {
		return false
	}
	for i := 1; i < len(name); i++ {
		if name[i] < '0' ||
			name[i] > 'z' ||
			name[i] > '9' && name[i] < 'A' ||
			name[i] > 'Z' && name[i] != '_' && name[i] < 'a' {
			return false
		}
	}
	return true
}

// PropName is the property name of semantic nodes and links.
//
// A valid property name consists of alphanumeric characters and
// underscores ('_'), begins with a lowercase letter,
// does not begin with "sln", and is up to 65535 bytes long.
type PropName struct {
	name string
}

// NewPropName returns a PropName whose value is propName.
//
// A valid property name consists of alphanumeric characters and
// underscores ('_'), begins with a lowercase letter,
// does not begin with "sln", and is up to 65535 bytes long.
// If propName is invalid, NewPropName reports a *InvalidPropNameError.
// (To test whether err is *InvalidPropNameError, use function errors.As.)
func NewPropName(propName string) (pn PropName, err error) {
	if IsValidPropNameString(propName) {
		pn.name = propName
	} else {
		err = errors.AutoWrap(NewInvalidPropNameError(propName))
	}
	return
}

// MustNewPropName is like NewPropName,
// but panic if the property name is invalid.
func MustNewPropName(propName string) PropName {
	pn, err := NewPropName(propName)
	if err != nil {
		panic(errors.AutoWrap(err))
	}
	return pn
}

// String returns the property name.
//
// If pn is invalid, String returns an empty string.
func (pn PropName) String() string {
	return pn.name
}

// IsValid reports whether pn is valid.
func (pn PropName) IsValid() bool {
	// Its constructor should guarantee that pn is valid if it is not zero.
	return pn.name != ""
}

// PropNameSet is a set of property names, all of which are valid PropName.
//
// If an invalid PropName is about to be put into this set,
// the corresponding method panics with a *InvalidPropNameError.
//
// To test whether the panic value is a *InvalidPropNameError,
// convert it to an error with type assertion,
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
type PropNameSet interface {
	set.Set[PropName]
}

// NewPropNameSet creates a new PropNameSet.
//
// The method Range of the set accesses property names in random order.
// The access order in two calls to Range may be different.
//
// capacity asks to allocate enough space to hold
// the specified number of property names.
// If capacity is negative, it is ignored.
func NewPropNameSet(capacity int) PropNameSet {
	return newValidSet(
		capacity,
		func(x PropName) bool {
			return x.IsValid()
		},
		func(x PropName) error {
			return NewInvalidPropNameError(x.String())
		},
	)
}

// mutExclPropNameSet is an implementation of interface PropNameSet.
//
// It can associate with one or more collections
// that have the method Remove(...PropName).
// When a property name is put into this set,
// mutExclPropNameSet removes the property name from these collections.
//
// The client must call its method init to initialize
// the mutExclPropNameSet before use.
type mutExclPropNameSet struct {
	s PropNameSet
	r []interface{ Remove(...PropName) }
}

// init initializes the mutExclPropNameSet
// with the specified capacity and collections.
//
// capacity asks to allocate enough space to hold
// the specified number of property names.
// If capacity is negative, it is ignored.
//
// collection is a list of collections associated with this set.
// When a property name is put into this set,
// mutExclPropNameSet removes the property name from these collections.
func (mepns *mutExclPropNameSet) init(capacity int,
	collection ...interface{ Remove(...PropName) }) {
	mepns.s = NewPropNameSet(capacity)
	if len(collection) > 0 {
		mepns.r = make([]interface{ Remove(...PropName) }, len(collection))
		copy(mepns.r, collection)
	}
}

func (mepns *mutExclPropNameSet) Len() int {
	mepns.checkInit()
	return mepns.s.Len()
}

// Range accesses the property names in the set.
// Each property name is accessed once.
// The access order may be random and may be different at each call.
//
// Its parameter handler is a function to deal with the property name
// in the set and report whether to continue to access the next property name.
func (mepns *mutExclPropNameSet) Range(handler func(x PropName) (cont bool)) {
	mepns.checkInit()
	mepns.s.Range(handler)
}

func (mepns *mutExclPropNameSet) Filter(filter func(x PropName) (keep bool)) {
	mepns.checkInit()
	mepns.s.Filter(filter)
}

func (mepns *mutExclPropNameSet) ContainsItem(x PropName) bool {
	mepns.checkInit()
	return mepns.s.ContainsItem(x)
}

func (mepns *mutExclPropNameSet) ContainsSet(s set.Set[PropName]) bool {
	mepns.checkInit()
	return mepns.s.ContainsSet(s)
}

func (mepns *mutExclPropNameSet) ContainsAny(
	c container.Container[PropName]) bool {
	mepns.checkInit()
	return mepns.s.ContainsAny(c)
}

func (mepns *mutExclPropNameSet) Add(x ...PropName) {
	mepns.checkInit()
	mepns.s.Add(x...)
	mepns.removeFromOthers(x...)
}

func (mepns *mutExclPropNameSet) Remove(x ...PropName) {
	mepns.checkInit()
	mepns.s.Remove(x...)
}

func (mepns *mutExclPropNameSet) Union(s set.Set[PropName]) {
	mepns.checkInit()
	if s == nil || s.Len() == 0 {
		return
	}
	mepns.s.Union(s)
	s.Range(func(x PropName) (cont bool) {
		mepns.removeFromOthers(x)
		return true
	})
}

func (mepns *mutExclPropNameSet) Intersect(s set.Set[PropName]) {
	mepns.checkInit()
	mepns.s.Intersect(s)
}

func (mepns *mutExclPropNameSet) Subtract(s set.Set[PropName]) {
	mepns.checkInit()
	mepns.s.Subtract(s)
}

func (mepns *mutExclPropNameSet) DisjunctiveUnion(s set.Set[PropName]) {
	mepns.checkInit()
	if s == nil || s.Len() == 0 {
		return
	}
	mepns.s.DisjunctiveUnion(s)
	s.Range(func(x PropName) (cont bool) {
		if mepns.s.ContainsItem(x) {
			mepns.removeFromOthers(x)
		}
		return true
	})
}

func (mepns *mutExclPropNameSet) Clear() {
	mepns.checkInit()
	mepns.s.Clear()
}

// checkInit checks whether mepns is initialized.
// If not, it panics.
func (mepns *mutExclPropNameSet) checkInit() {
	if mepns.s == nil {
		panic(errors.AutoMsgCustom("not initialized before use", -1, 1))
	}
}

// removeFromOthers removes name from collections in mepns.r.
func (mepns *mutExclPropNameSet) removeFromOthers(name ...PropName) {
	if len(name) > 0 {
		for _, r := range mepns.r {
			r.Remove(name...)
		}
	}
}
