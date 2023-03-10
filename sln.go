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
// The CURD operations after Close report ErrSLNClosed.
// (To test whether an error is ErrSLNClosed, use function errors.Is.)
// The successive calls to Close do nothing
// but block until the SLN is closed or any error occurs during closing.
type SLN interface {
	inout.Closer

	// NumNodeType returns the number of node types and any error encountered.
	NumNodeType(ctx context.Context) (n int, err error)

	// NumLinkType returns the number of link types and any error encountered.
	NumLinkType(ctx context.Context) (n int, err error)

	// NumNode returns the number of nodes that satisfy
	// the specified conditions and any error encountered.
	NumNode(ctx context.Context, cond NodeMatchCond) (n int, err error)

	// NumLink returns the number of links that satisfy
	// the specified conditions and any error encountered.
	NumLink(ctx context.Context, cond LinkMatchCond) (n int, err error)

	// GetNodeTypes returns all node types in this SLN.
	GetNodeTypes(ctx context.Context) (types []Type, err error)

	// GetLinkTypes returns all link types in this SLN.
	GetLinkTypes(ctx context.Context) (types []Type, err error)

	// GetNodeByID returns the node with the specified ID
	// and any error encountered.
	//
	// GetNodeByID reports a *NodeNotExistError if the node does not exist.
	// (To test whether err is *NodeNotExistError, use function errors.As.)
	//
	// propTypes specify the types of properties on the node.
	// The properties not in propTypes are discarded.
	//
	// GetNodeByID reports a *PropTypeError if any property
	// does not match its specified type.
	// (To test whether err is *PropTypeError, use function errors.As.)
	GetNodeByID(ctx context.Context, id ID, propTypes PropTypeMap) (node *Node, err error)

	// GetLinkByID returns the link with the specified ID
	// and any error encountered.
	//
	// GetLinkByID reports a *LinkNotExistError if the link does not exist.
	// (To test whether err is *LinkNotExistError, use function errors.As.)
	//
	// propTypes specify the types of properties on the link.
	// The properties not in propTypes are discarded.
	//
	// GetLinkByID reports a *PropTypeError if any property
	// does not match its specified type.
	// (To test whether err is *PropTypeError, use function errors.As.)
	GetLinkByID(ctx context.Context, id ID, propTypes PropTypeMap) (link *Link, err error)

	// GetAllNodes returns all nodes that satisfy the specified conditions
	// and any error encountered.
	//
	// propTypes specify the types of properties on the node.
	// The properties not in propTypes are discarded.
	//
	// GetAllNodes reports a *PropTypeError if any property
	// does not match its specified type.
	// (To test whether err is *PropTypeError, use function errors.As.)
	GetAllNodes(ctx context.Context, propTypes PropTypeMap, cond NodeMatchCond) (nodes []*Node, err error)

	// GetAllLinks returns all links that satisfy the specified conditions
	// and any error encountered.
	//
	// propTypes specify the types of properties on the link.
	// The properties not in propTypes are discarded.
	//
	// GetAllLinks reports a *PropTypeError if any property
	// does not match its specified type.
	// (To test whether err is *PropTypeError, use function errors.As.)
	GetAllLinks(ctx context.Context, propTypes PropTypeMap, cond LinkMatchCond) (links []*Link, err error)

	// CreateNode creates a new node with the specified node type t.
	//
	// props are initial properties on the new node.
	//
	// CreateNode reports a *InvalidTypeError if t is invalid.
	// (To test whether err is *InvalidTypeError, use function errors.As.)
	CreateNode(ctx context.Context, t Type, props PropMap) (node *Node, err error)

	// CreateLink creates a new link with the specified link type t,
	// starting from the node with ID "from" and
	// pointing to the node with ID "to".
	//
	// props are initial properties on the new link.
	//
	// CreateLink reports a *InvalidTypeError if t is invalid.
	// (To test whether err is *InvalidTypeError, use function errors.As.)
	//
	// CreateLink reports a *NodeNotExistError if from or to does not exist.
	// (To test whether err is *NodeNotExistError, use function errors.As.)
	CreateLink(ctx context.Context, t Type, from, to ID, props PropMap) (link *Link, err error)

	// RemoveNodeByID removes the node with the specified ID
	// and all associated links.
	//
	// It returns nil error if there is no such node or id is invalid.
	RemoveNodeByID(ctx context.Context, id ID) error

	// RemoveLinkByID removes the link with the specified ID.
	//
	// It returns nil error if there is no such link or id is invalid.
	RemoveLinkByID(ctx context.Context, id ID) error

	// SetNodeProperties sets the properties on the node
	// that has the specified ID to the specified properties.
	//
	// It removes all properties on the node if props are nil or empty.
	//
	// It returns the node updated and any error encountered.
	SetNodeProperties(ctx context.Context, id ID, props PropMap) (node *Node, err error)

	// SetLinkProperties sets the properties on the link
	// that has the specified ID to the specified properties.
	//
	// It removes all properties on the link if props are nil or empty.
	//
	// It returns the link updated and any error encountered.
	SetLinkProperties(ctx context.Context, id ID, props PropMap) (link *Link, err error)

	// MutateNodeProperties mutates the properties on the node
	// that has the specified ID.
	//
	// It returns the node updated and any error encountered.
	MutateNodeProperties(ctx context.Context, id ID, pma PropMutateArg) (node *Node, err error)

	// MutateLinkProperties mutates the properties on the link
	// that has the specified ID.
	//
	// It returns the link updated and any error encountered.
	MutateLinkProperties(ctx context.Context, id ID, pma PropMutateArg) (link *Link, err error)
}

// NL consists of the common fields of Node and Link.
type NL struct {
	SLN   SLN     // The Semantic Link Network to which this node or link belongs.
	ID    ID      // The ID of this node or link.
	Type  Type    // The type of this node or link.
	Props PropMap // The properties on this node or link.
}

// Node records the information of a semantic node.
type Node struct {
	NL
}

// Link records the information of a semantic link.
type Link struct {
	NL
	From *Node // The node from which this link starts.
	To   *Node // The node to which this link points.
}
