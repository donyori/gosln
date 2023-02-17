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
	"context"

	"github.com/donyori/gogo/inout"
)

// SLN contains basic CRUD (create, read, update, and delete)
// operations on the Semantic Link Network.
//
// It is safe for concurrency.
//
// For each CRUD operation, it supports the client in using
// context.Context to set a deadline or a cancellation signal.
// If the deadline and cancellation are not required,
// the client should pass a context.Background() instead of nil.
//
// Its method Close marks the SLN as unusable and releases the resource.
// Close waits for the in-flight CURD operations rather than interrupting them.
// The CURD operations after Close will report ErrSLNClosed.
// (To test whether an error is ErrSLNClosed, use function errors.Is.)
// The successive calls to Close will do nothing
// but block until the SLN is closed or any error occurs during closing.
type SLN interface {
	inout.Closer

	// NumNodeType returns the number of node types and any error encountered.
	NumNodeType(ctx context.Context) (n int, err error)

	// NumLinkType returns the number of link types and any error encountered.
	NumLinkType(ctx context.Context) (n int, err error)

	// NumNode returns the number of nodes with the specified type
	// and any error encountered.
	// In particular, if the type is invalid (such as zero-value),
	// it returns the number of all nodes.
	NumNode(ctx context.Context, t Type) (n int, err error)

	// NumLink returns the number of links with the specified type
	// and any error encountered.
	// In particular, if the type is invalid (such as zero-value),
	// it returns the number of all links.
	NumLink(ctx context.Context, t Type) (n int, err error)

	// GetNodeTypes returns all node types in this SLN.
	GetNodeTypes(ctx context.Context) (types []Type, err error)

	// GetLinkTypes returns all link types in this SLN.
	GetLinkTypes(ctx context.Context) (types []Type, err error)

	// GetNode returns the node with the specified ID
	// and any error encountered.
	GetNode(ctx context.Context, id ID) (node *Node, err error)

	// GetLink returns the link with the specified ID
	// and any error encountered.
	GetLink(ctx context.Context, id ID) (link *Link, err error)

	// GetAllNodes returns all nodes with the specified type
	// and any error encountered.
	// In particular, if the type is invalid (such as zero-value),
	// it returns all nodes in this SLN.
	GetAllNodes(ctx context.Context, t Type) (nodes []*Node, err error)

	// GetAllLinks returns all links with the specified type
	// and any error encountered.
	// In particular, if the type is invalid (such as zero-value),
	// it returns all links in this SLN.
	GetAllLinks(ctx context.Context, t Type) (links []*Link, err error)

	// CreateNode creates a new node with the specified node type t.
	//
	// prop is a set of initial properties of the new node.
	//
	// It reports a *InvalidTypeError if t is invalid.
	// (To test whether err is *InvalidTypeError, use function errors.As.)
	CreateNode(ctx context.Context, t Type, prop *PropertyMap) (node *Node, err error)

	// CreateLink creates a new link with the specified link type t,
	// starting from the node with ID from and pointing to the node with ID to.
	//
	// prop is a set of initial properties of the new link.
	//
	// It reports a *InvalidTypeError if t is invalid.
	// (To test whether err is *InvalidTypeError, use function errors.As.)
	//
	// It reports a *NodeNotExistError if from or to does not exist.
	// (To test whether err is *NodeNotExistError, use function errors.As.)
	CreateLink(ctx context.Context, t Type, from, to ID, prop *PropertyMap) (link *Link, err error)

	// RemoveNode removes the node with the specified ID
	// and all associated links.
	// It returns the properties of that node and any error encountered.
	//
	// It returns nil PropertyMap and nil error
	// if there is no such node or id is invalid.
	RemoveNode(ctx context.Context, id ID) (prop *PropertyMap, err error)

	// RemoveLink removes the link with the specified ID.
	// It returns the properties of that link and any error encountered.
	//
	// It returns nil PropertyMap and nil error
	// if there is no such link or id is invalid.
	RemoveLink(ctx context.Context, id ID) (prop *PropertyMap, err error)

	// UpdateNodeProperty updates the properties of
	// the node with the specified ID.
	// It returns the node updated and any error encountered.
	UpdateNodeProperty(ctx context.Context, id ID, prop *PropertyMap) (node *Node, err error)

	// UpdateLinkProperty updates the properties of
	// the link with the specified ID.
	// It returns the link updated and any error encountered.
	UpdateLinkProperty(ctx context.Context, id ID, prop *PropertyMap) (link *Link, err error)
}

// NL consists of the common fields of Node and Link.
type NL struct {
	SLN  SLN          // The Semantic Link Network to which this node or link belongs.
	ID   ID           // The ID of this node or link.
	Type Type         // The type of this node or link.
	Prop *PropertyMap // The properties of this node or link.
}

// Node records the information of a semantic node.
type Node struct {
	NL
	Outgoing *IDSet // The IDs of its outgoing links.
	Incoming *IDSet // The IDs of its incoming links.
}

// Link records the information of a semantic link.
type Link struct {
	NL
	From ID // The ID of the node where this link starts.
	To   ID // The ID of the node to which this link points.
}
