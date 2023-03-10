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
	"reflect"
	"strconv"
	"strings"

	"github.com/donyori/gogo/errors"
	"github.com/donyori/gogo/inout"
)

// ErrSLNClosed is an error indicating that the SLN is already closed.
//
// The client should use errors.Is to test whether an error is ErrSLNClosed.
var ErrSLNClosed = errors.AutoWrapCustom(
	inout.NewClosedError("SLN", nil),
	errors.PrependFullPkgName,
	0,
	nil,
)

// InvalidTypeError is an error indicating that the type is invalid.
type InvalidTypeError struct {
	t string // The type, as a string.
}

var _ error = (*InvalidTypeError)(nil)

// NewInvalidTypeError creates a new InvalidTypeError
// with the specified type t.
func NewInvalidTypeError(t string) *InvalidTypeError {
	return &InvalidTypeError{t: t}
}

// Type returns the type recorded in e, as a string.
//
// If e is nil, it returns "<nil>".
func (e *InvalidTypeError) Type() string {
	if e == nil {
		return "<nil>"
	}
	return e.t
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *InvalidTypeError>".
func (e *InvalidTypeError) Error() string {
	if e == nil {
		return "<nil *InvalidTypeError>"
	}
	return "type " + strconv.Quote(e.t) + " is invalid; " +
		"a valid type consists of alphanumeric characters and underscores ('_'), " +
		`begins with an uppercase letter, does not begin with "SLN", ` +
		"and is up to 65535 bytes long."
}

// InvalidIDError is an error indicating that the ID is invalid.
type InvalidIDError struct {
	id ID
}

var _ error = (*InvalidIDError)(nil)

// NewInvalidIDError creates a new InvalidIDError with the specified ID.
func NewInvalidIDError(id ID) *InvalidIDError {
	return &InvalidIDError{id: id}
}

// ID returns the ID recorded in e.
//
// If e is nil, it returns a zero-value ID.
func (e *InvalidIDError) ID() ID {
	if e == nil {
		return ID{}
	}
	return e.id
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *InvalidIDError>".
func (e *InvalidIDError) Error() string {
	if e == nil {
		return "<nil *InvalidTypeError>"
	}
	return "ID " + strconv.Quote(e.id.String()) + " is invalid"
}

// InvalidPropNameError is an error indicating that
// the property name is invalid.
type InvalidPropNameError struct {
	name string // The property name, as a string.
}

var _ error = (*InvalidPropNameError)(nil)

// NewInvalidPropNameError creates a new InvalidPropNameError
// with the specified property name.
func NewInvalidPropNameError(propName string) *InvalidPropNameError {
	return &InvalidPropNameError{name: propName}
}

// PropName returns the property name recorded in e, as a string.
//
// If e is nil, it returns "<nil>".
func (e *InvalidPropNameError) PropName() string {
	if e == nil {
		return "<nil>"
	}
	return e.name
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *InvalidPropNameError>".
func (e *InvalidPropNameError) Error() string {
	if e == nil {
		return "<nil *InvalidPropNameError>"
	}
	return "property name " + strconv.Quote(e.name) + " is invalid; " +
		"a valid property name consists of alphanumeric characters and underscores ('_'), " +
		`begins with a lowercase letter, does not begin with "sln", ` +
		"and is up to 65535 bytes long."
}

// InvalidPropTypeError is an error indicating that
// the property type is invalid.
type InvalidPropTypeError struct {
	t PropType // The property type.
}

var _ error = (*InvalidPropTypeError)(nil)

// NewInvalidPropTypeError creates a new InvalidPropTypeError
// with the specified property type.
func NewInvalidPropTypeError(propType PropType) *InvalidPropTypeError {
	return &InvalidPropTypeError{t: propType}
}

// PropType returns the property type recorded in e.
//
// If e is nil, it returns 0.
func (e *InvalidPropTypeError) PropType() PropType {
	if e == nil {
		return 0
	}
	return e.t
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *InvalidPropTypeError>".
func (e *InvalidPropTypeError) Error() string {
	if e == nil {
		return "<nil *InvalidPropTypeError>"
	}
	return "property type " + e.t.String() + " is invalid"
}

// InvalidPropValueError is an error indicating that
// the property value is invalid.
type InvalidPropValueError struct {
	value any // The property value.
}

var _ error = (*InvalidPropValueError)(nil)

// NewInvalidPropValueError creates a new InvalidPropValueError
// with the specified property value.
func NewInvalidPropValueError(propValue any) *InvalidPropValueError {
	return &InvalidPropValueError{value: propValue}
}

// PropValue returns the property value recorded in e.
//
// If e is nil, it returns nil.
func (e *InvalidPropValueError) PropValue() any {
	if e == nil {
		return nil
	}
	return e.value
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *InvalidPropValueError>".
func (e *InvalidPropValueError) Error() string {
	if e == nil {
		return "<nil *InvalidPropValueError>"
	}
	var b strings.Builder
	b.WriteString("property value (type: ")
	_, _ = fmt.Fprintf(&b, "%#v", e.value) // ignore error as it is always nil
	b.WriteString(") is invalid; the type of valid property value must be one of ")
	for i := PropType(1); i.IsValid(); i++ {
		if i > 1 {
			b.WriteString(", ")
		}
		b.WriteString(i.String())
	}
	return b.String()
}

// PropNotExistError is an error indicating that
// the property with the specified name does not exist.
type PropNotExistError struct {
	name PropName // The property name.
}

var _ error = (*PropNotExistError)(nil)

// NewPropNotExistError creates a new PropNotExistError
// with the specified property name.
func NewPropNotExistError(propName PropName) *PropNotExistError {
	return &PropNotExistError{name: propName}
}

// PropName returns the property name recorded in e.
//
// If e is nil, it returns a zero-value PropName.
func (e *PropNotExistError) PropName() PropName {
	if e == nil {
		return PropName{}
	}
	return e.name
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *PropNotExistError>".
func (e *PropNotExistError) Error() string {
	if e == nil {
		return "<nil *PropNotExistError>"
	}
	name := e.name.String()
	if name == "" {
		name = "property"
	}
	return name + " does not exist"
}

// PropTypeError is an error indicating that the property type is wrong.
//
// It records the property name, value, and expected type.
type PropTypeError struct {
	name     PropName     // The property name.
	value    any          // The property value.
	wantType reflect.Type // The expected type.
}

var _ error = (*PropTypeError)(nil)

// NewPropTypeError creates a new PropTypeError with
// the specified property name, value, and expected type.
func NewPropTypeError(
	propName PropName,
	propValue any,
	wantType reflect.Type,
) *PropTypeError {
	return &PropTypeError{
		name:     propName,
		value:    propValue,
		wantType: wantType,
	}
}

// PropName returns the property name recorded in e.
//
// If e is nil, it returns a zero-value PropName.
func (e *PropTypeError) PropName() PropName {
	if e == nil {
		return PropName{}
	}
	return e.name
}

// PropValue returns the property value recorded in e.
//
// If e is nil, it returns nil.
func (e *PropTypeError) PropValue() any {
	if e == nil {
		return nil
	}
	return e.value
}

// WantType returns the expected type recorded in e.
//
// If e is nil, it returns nil.
func (e *PropTypeError) WantType() reflect.Type {
	if e == nil {
		return nil
	}
	return e.wantType
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *PropTypeError>".
func (e *PropTypeError) Error() string {
	if e == nil {
		return "<nil *PropTypeError>"
	}
	var b strings.Builder
	b.WriteString(e.name.String())
	if b.Len() == 0 {
		b.WriteString("property")
	}
	b.WriteString(" has wrong type ")
	b.WriteString(reflect.TypeOf(e.value).String())
	if e.wantType != nil {
		b.WriteString("; want ")
		b.WriteString(e.wantType.String())
	}
	return b.String()
}

// NodeNotExistError is an error indicating that
// the node with the specified ID does not exist.
type NodeNotExistError struct {
	id ID // The node ID.
}

var _ error = (*NodeNotExistError)(nil)

// NewNodeNotExistError creates a new NodeNotExistError
// with the specified node ID.
func NewNodeNotExistError(nodeID ID) *NodeNotExistError {
	return &NodeNotExistError{id: nodeID}
}

// NodeID returns the node ID recorded in e.
//
// If e is nil, it returns a zero-value ID (invalid).
func (e *NodeNotExistError) NodeID() ID {
	if e == nil {
		return ID{}
	}
	return e.id
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *NodeNotExistError>".
func (e *NodeNotExistError) Error() string {
	if e == nil {
		return "<nil *NodeNotExistError>"
	}
	return "node " + strconv.Quote(e.id.String()) + " does not exist"
}

// LinkNotExistError is an error indicating that
// the link with the specified ID does not exist.
type LinkNotExistError struct {
	id ID // The link ID.
}

var _ error = (*LinkNotExistError)(nil)

// NewLinkNotExistError creates a new LinkNotExistError
// with the specified link ID.
func NewLinkNotExistError(linkID ID) *LinkNotExistError {
	return &LinkNotExistError{id: linkID}
}

// LinkID returns the link ID recorded in e.
//
// If e is nil, it returns a zero-value ID (invalid).
func (e *LinkNotExistError) LinkID() ID {
	if e == nil {
		return ID{}
	}
	return e.id
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *LinkNotExistError>".
func (e *LinkNotExistError) Error() string {
	if e == nil {
		return "<nil *LinkNotExistError>"
	}
	return "link " + strconv.Quote(e.id.String()) + " does not exist"
}
