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
	"strconv"

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
		"begins with an uppercase letter, and is up to 65535 bytes long."
}

// InvalidPropertyNameError is an error indicating that
// the property name is invalid.
type InvalidPropertyNameError struct {
	name string // The property name.
}

var _ error = (*InvalidPropertyNameError)(nil)

// NewInvalidPropertyNameError creates a new InvalidPropertyNameError
// with the specified property name.
func NewInvalidPropertyNameError(propertyName string) *InvalidPropertyNameError {
	return &InvalidPropertyNameError{name: propertyName}
}

// PropertyName returns the property name recorded in e.
//
// If e is nil, it returns "<nil>".
func (e *InvalidPropertyNameError) PropertyName() string {
	if e == nil {
		return "<nil>"
	}
	return e.name
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *InvalidPropertyNameError>".
func (e *InvalidPropertyNameError) Error() string {
	if e == nil {
		return "<nil *InvalidPropertyNameError>"
	}
	return "property name " + strconv.Quote(e.name) + " is invalid; " +
		"a valid property name consists of alphanumeric characters and underscores ('_'), " +
		"begins with a lowercase letter, and is up to 65535 bytes long."
}

// PropertyNotExistError is an error indicating that
// the property with the specified name does not exist.
type PropertyNotExistError struct {
	name string // The property name.
}

var _ error = (*PropertyNotExistError)(nil)

// NewPropertyNotExistError creates a new PropertyNotExistError
// with the specified property name.
func NewPropertyNotExistError(propertyName string) *PropertyNotExistError {
	return &PropertyNotExistError{name: propertyName}
}

// PropertyName returns the property name recorded in e.
//
// If e is nil, it returns "<nil>".
func (e *PropertyNotExistError) PropertyName() string {
	if e == nil {
		return "<nil>"
	}
	return e.name
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *PropertyNotExistError>".
func (e *PropertyNotExistError) Error() string {
	if e == nil {
		return "<nil *PropertyNotExistError>"
	}
	name := e.name
	if name == "" {
		name = "property"
	}
	return e.name + " does not exist"
}

// PropertyTypeError is an error indicating that the property type is wrong.
//
// It records the property name, value, and the name of the expected type.
type PropertyTypeError struct {
	name     string // The property name.
	value    any    // The property value.
	wantType string // The name of the expected type.
}

var _ error = (*PropertyTypeError)(nil)

// NewPropertyTypeError creates a new PropertyTypeError with the specified
// property name, value, and the name of the expected type.
func NewPropertyTypeError(
	propertyName string,
	propertyValue any,
	wantType string,
) *PropertyTypeError {
	return &PropertyTypeError{
		name:     propertyName,
		value:    propertyValue,
		wantType: wantType,
	}
}

// PropertyName returns the property name recorded in e.
//
// If e is nil, it returns "<nil>".
func (e *PropertyTypeError) PropertyName() string {
	if e == nil {
		return "<nil>"
	}
	return e.name
}

// PropertyValue returns the property value recorded in e.
//
// If e is nil, it returns nil.
func (e *PropertyTypeError) PropertyValue() any {
	if e == nil {
		return nil
	}
	return e.value
}

// WantType returns the name of the expected type recorded in e.
//
// If e is nil, it returns "<nil>".
func (e *PropertyTypeError) WantType() string {
	if e == nil {
		return "<nil>"
	}
	return e.wantType
}

// Error returns the error message.
//
// If e is nil, it returns "<nil *PropertyTypeError>".
func (e *PropertyTypeError) Error() string {
	if e == nil {
		return "<nil *PropertyTypeError>"
	}
	name := e.name
	if name == "" {
		name = "property"
	}
	s := fmt.Sprintf("%s has wrong type %T", name, e.value)
	if e.wantType != "" {
		s += "; want " + e.wantType
	}
	return s
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
	return fmt.Sprintf("node %q does not exist", e.id)
}
