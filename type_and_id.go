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
	"regexp"
	"strings"

	"github.com/donyori/gogo/errors"
)

// encode64Table is a character table for encoding
// the serial number (int64) to a valid suffix of ID.
const encode64Table = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

// typePattern is the regular expression pattern for Type.
var typePattern = regexp.MustCompile("^[A-Z][0-9A-Z_a-z]{0,65534}$")

// Type is the type of the semantic node and link.
//
// A valid type consists of alphanumeric characters and underscores ('_'),
// begins with an uppercase letter, and is up to 65535 bytes long.
type Type struct {
	t string
}

// NewType returns a Type whose value is t.
//
// A valid type consists of alphanumeric characters and underscores ('_'),
// begins with an uppercase letter, and is up to 65535 bytes long.
// If t is invalid, NewType will report a *InvalidTypeError.
// (To test whether err is *InvalidTypeError, use function errors.As.)
func NewType(t string) (typ Type, err error) {
	if typePattern.MatchString(t) {
		typ = Type{t: t}
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
// If t is invalid, String will return an empty string.
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

// NewID returns an ID corresponding to the type t, with the serial number i.
//
// If t is invalid (such as zero-value), NewID will return a zero-value ID.
//
// If i is negative, NewID will panic.
func NewID(t Type, i int64) ID {
	if i < 0 {
		panic(errors.AutoMsg(fmt.Sprintf("the number i (%d) is negative", i)))
	}
	if !t.IsValid() {
		return ID{}
	}
	var b strings.Builder
	b.Grow(11)
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

// String formats id to a string, as follows:
//
//	<Type>#<UniqueSuffix>
//
// where <Type> is the type corresponding to id,
// and <UniqueSuffix> is a suffix that is unique
// across the IDs for the same type.
//
// If id is invalid, String will return an empty string.
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

// IDSet is a set of IDs.
//
// A zero-value IDSet is ready to use.
type IDSet struct {
	m            map[string]map[string]struct{}
	containsZero bool
}

// Len returns the number of IDs in the set.
func (s *IDSet) Len() int {
	if s == nil {
		return 0
	}
	var n int
	if s.containsZero {
		n = 1
	}
	for _, sub := range s.m {
		n += len(sub)
	}
	return n
}

// LenType returns the number of IDs corresponding to the type t in the set.
func (s *IDSet) LenType(t Type) int {
	switch {
	case s == nil:
	case t.t == "":
		if s.containsZero {
			return 1
		}
	case len(s.m) == 0:
	default:
		return len(s.m[t.t])
	}
	return 0
}

// NumType returns the number of types corresponding to the IDs in the set.
func (s *IDSet) NumType() int {
	if s == nil {
		return 0
	}
	var x int
	if s.containsZero {
		x = 1
	}
	return len(s.m) + x
}

// ContainsType reports whether there is an ID
// corresponding to the type t in the set.
func (s *IDSet) ContainsType(t Type) bool {
	switch {
	case s == nil:
	case t.t == "":
		return s.containsZero
	case len(s.m) == 0:
	default:
		return len(s.m[t.t]) > 0
	}
	return false
}

// ContainsID reports whether id is in the set.
func (s *IDSet) ContainsID(id ID) bool {
	switch {
	case s == nil:
	case id.t == "":
		return s.containsZero
	case len(s.m) == 0:
	default:
		sub := s.m[id.t]
		if sub == nil {
			return false
		}
		_, ok := sub[id.s]
		return ok
	}
	return false
}

// Range accesses the IDs in the set.
// Each ID will be accessed once.
// The order of the access is random.
//
// Its parameter handler is a function to deal with an ID in the set
// and report whether to continue to access the next ID.
func (s *IDSet) Range(handler func(id ID) (cont bool)) {
	if s == nil || s.containsZero && !handler(ID{}) {
		return
	}
	for t, sub := range s.m {
		for suffix := range sub {
			if !handler(ID{t: t, s: suffix}) {
				return
			}
		}
	}
}

// RangeType accesses the IDs corresponding to the type t in the set.
// Each ID will be accessed once. The order of the access is random.
//
// Its parameter handler is a function to deal with an ID in the set
// and report whether to continue to access the next ID.
func (s *IDSet) RangeType(t Type, handler func(id ID) (cont bool)) {
	switch {
	case s == nil:
	case t.t == "":
		if s.containsZero {
			handler(ID{})
		}
	case len(s.m) == 0:
	default:
		for suffix := range s.m[t.t] {
			if !handler(ID{t: t.t, s: suffix}) {
				return
			}
		}
	}
}

// Add adds id to the set.
func (s *IDSet) Add(id ...ID) {
	if s == nil {
		panic(errors.AutoMsg("*IDSet is nil"))
	}
	if len(id) == 0 {
		return
	}
	for _, x := range id {
		if x.t == "" {
			s.containsZero = true
			continue
		}
		if s.m == nil {
			s.m = make(map[string]map[string]struct{})
		}
		sub := s.m[x.t]
		if sub == nil {
			sub = make(map[string]struct{})
			s.m[x.t] = sub
		}
		sub[x.s] = struct{}{}
	}
}

// Remove removes id from the set.
// It does nothing for those that are not in the set.
func (s *IDSet) Remove(id ...ID) {
	if s == nil || len(id) == 0 {
		return
	}
	for _, x := range id {
		if x.t == "" {
			s.containsZero = false
		} else if len(s.m) > 0 {
			sub := s.m[x.t]
			if sub != nil {
				delete(sub, x.s)
				if len(sub) == 0 {
					delete(s.m, x.t)
				}
			}
		}
	}
}

// Clear removes all IDs in the set.
func (s *IDSet) Clear() {
	s.m = nil
	s.containsZero = false
}
