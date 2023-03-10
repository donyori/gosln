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

import "github.com/donyori/gogo/container/mapping"

// PropMatchClause is a conjunction of conditions to
// match properties on a semantic node or link.
//
// A set of properties satisfies the PropMatchClause
// if it satisfies all the conditions in this PropMatchClause.
//
// PropMatchClause consists of three components:
//   - Equal: a PropMap holding the properties that must be equal to the target properties.
//   - Present: a PropNameSet holding the names of the properties that must exist.
//   - Absent: a PropNameSet holding the names of the properties that must not exist.
//
// These components are mutually exclusive:
// when a property is put into one component, it is removed from the others.
type PropMatchClause interface {
	// Equal returns a PropMap with properties
	// that must be equal to the target properties.
	//
	// The PropMap is always non-nil, but may be empty.
	Equal() PropMap

	// Present returns a PropNameSet with
	// names of the properties that must exist.
	//
	// The PropNameSet is always non-nil, but may be empty.
	Present() PropNameSet

	// Absent returns a PropNameSet with
	// names of the properties that must not exist.
	//
	// The PropNameSet is always non-nil, but may be empty.
	Absent() PropNameSet

	// Match reports whether props satisfy this PropMatchClause.
	Match(props PropMap) bool
}

// propMatchClauseImpl is an implementation of interface PropMatchClause.
type propMatchClauseImpl struct {
	equal   *mutExclPropMap     // Properties that must be equal to the target properties.
	present *mutExclPropNameSet // Names of the properties that must exist.
	absent  *mutExclPropNameSet // Names of the properties that must not exist.
}

// NewPropMatchClause creates a new PropMatchClause.
//
// eqCap, presentCap, and absentCap ask to allocate enough space to hold
// the specified number of items in its Equal, Present, and Absent components,
// respectively.
// If eqCap is negative, it is ignored, as are presentCap and absentCap.
func NewPropMatchClause(eqCap, presentCap, absentCap int) PropMatchClause {
	pmc := &propMatchClauseImpl{
		equal:   new(mutExclPropMap),
		present: new(mutExclPropNameSet),
		absent:  new(mutExclPropNameSet),
	}
	pmc.equal.init(eqCap, pmc.present, pmc.absent)
	pmc.present.init(presentCap, pmc.equal, pmc.absent)
	pmc.absent.init(absentCap, pmc.equal, pmc.present)
	return pmc
}

func (pmc *propMatchClauseImpl) Equal() PropMap {
	return pmc.equal
}

func (pmc *propMatchClauseImpl) Present() PropNameSet {
	return pmc.present
}

func (pmc *propMatchClauseImpl) Absent() PropNameSet {
	return pmc.absent
}

func (pmc *propMatchClauseImpl) Match(props PropMap) bool {
	if props == nil {
		return pmc.equal.Len() == 0 && pmc.present.Len() == 0
	}
	var ok bool
	pmc.equal.Range(func(x mapping.Entry[PropName, any]) (cont bool) {
		var value any
		value, ok = props.Get(x.Key)
		ok = ok && value == x.Value
		return ok
	})
	if !ok {
		return false
	}
	pmc.present.Range(func(x PropName) (cont bool) {
		_, ok = props.Get(x)
		return ok
	})
	if !ok {
		return false
	}
	pmc.absent.Range(func(x PropName) (cont bool) {
		_, ok = props.Get(x)
		return !ok
	})
	return !ok
}

// PropMatchCond is a disjunction of the clauses of type PropMatchClause
// to match properties on a semantic node or link.
//
// Any nil PropMatchClause in the PropMatchCond is ignored.
//
// A set of properties satisfies the PropMatchCond
// if it satisfies any of these clauses.
//
// In particular, a nil PropMatchCond matches any properties
// (including nil PropMap).
// A non-nil but empty PropMatchCond matches nothing.
type PropMatchCond []PropMatchClause

// Match reports whether props satisfy this PropMatchCond.
func (cond PropMatchCond) Match(props PropMap) bool {
	if cond == nil {
		return true
	}
	for _, pmc := range cond {
		if pmc != nil && pmc.Match(props) {
			return true
		}
	}
	return false
}

// NLMatchClause is a conjunction of conditions to
// match a semantic node or link.
//
// A semantic node or link satisfies the NLMatchClause
// if it satisfies all the conditions in this NLMatchClause.
//
// NLMatchClause can specify the ID, type,
// and properties on the semantic node or link.
type NLMatchClause interface {
	// GetID returns the specified ID of the semantic node or link.
	//
	// If no ID is specified, it returns a zero-value ID.
	GetID() ID

	// SetID specifies the ID of the semantic node or link.
	//
	// If id is invalid, it considers the ID unspecified.
	SetID(id ID)

	// SetIDAndClearOtherConds specifies the ID of the semantic node or link
	// and removes other match conditions.
	//
	// If id is invalid, it considers the ID unspecified.
	SetIDAndClearOtherConds(id ID)

	// GetType returns the specified type of the semantic node or link.
	//
	// If no type is specified, it returns a zero-value Type.
	GetType() Type

	// SetType specifies the type of the semantic node or link.
	//
	// If t is invalid, it considers the type unspecified.
	SetType(t Type)

	// GetPropMatchClause returns the match conditions for
	// properties on the semantic node or link.
	//
	// If there is no limit on the properties, it returns nil.
	GetPropMatchClause() PropMatchClause

	// SetPropMatchClause specifies the match conditions for
	// properties on the semantic node or link.
	//
	// If pmc is nil, it considers no limit on the properties.
	SetPropMatchClause(pmc PropMatchClause)
}

// nlMatchClauseImpl implements interface NLMatchClause,
// except for the method SetIDAndClearOtherConds.
type nlMatchClauseImpl struct {
	id  ID              // The specified ID, zero value for unspecified.
	t   Type            // The specified type, zero value for unspecified.
	pmc PropMatchClause // Match conditions for properties on the semantic node or link.
}

func (nlmc *nlMatchClauseImpl) GetID() ID {
	return nlmc.id
}

func (nlmc *nlMatchClauseImpl) SetID(id ID) {
	if id.IsValid() {
		nlmc.id = id
	} else {
		nlmc.id = ID{}
	}
}

func (nlmc *nlMatchClauseImpl) GetType() Type {
	return nlmc.t
}

func (nlmc *nlMatchClauseImpl) SetType(t Type) {
	if t.IsValid() {
		nlmc.t = t
	} else {
		nlmc.t = Type{}
	}
}

func (nlmc *nlMatchClauseImpl) GetPropMatchClause() PropMatchClause {
	return nlmc.pmc
}

func (nlmc *nlMatchClauseImpl) SetPropMatchClause(pmc PropMatchClause) {
	nlmc.pmc = pmc
}

// NodeMatchClause is a conjunction of conditions to match a semantic node.
//
// A semantic node satisfies the NodeMatchClause
// if it satisfies all the conditions in this NodeMatchClause.
//
// NodeMatchClause can specify the node ID, node type,
// and properties on the node.
type NodeMatchClause interface {
	NLMatchClause

	// Match reports whether the semantic node satisfies this NodeMatchClause.
	Match(node *Node) bool
}

// nodeMatchClauseImpl is an implementation of interface NodeMatchClause.
type nodeMatchClauseImpl struct {
	nlMatchClauseImpl
}

// NewNodeMatchClause creates a new NodeMatchClause.
func NewNodeMatchClause() NodeMatchClause {
	return new(nodeMatchClauseImpl)
}

func (nmc *nodeMatchClauseImpl) SetIDAndClearOtherConds(id ID) {
	nmc.SetID(id)
	nmc.t, nmc.pmc = Type{}, nil
}

func (nmc *nodeMatchClauseImpl) Match(node *Node) bool {
	switch {
	case node == nil:
	case nmc.id.IsValid() && node.ID != nmc.id:
	case nmc.t.IsValid() && node.Type != nmc.t:
	case nmc.pmc != nil && !nmc.pmc.Match(node.Props):
	default:
		return true
	}
	return false
}

// NodeMatchCond is a disjunction of the clauses of type NodeMatchClause
// to match a semantic node.
//
// Any nil NodeMatchClause in the NodeMatchCond is ignored.
//
// A semantic node satisfies the NodeMatchCond
// if it satisfies any of these clauses.
//
// In particular, a nil NodeMatchCond matches any semantic node (including nil).
// A non-nil but empty NodeMatchCond matches nothing.
type NodeMatchCond []NodeMatchClause

// Match reports whether the semantic node satisfies this NodeMatchCond.
func (cond NodeMatchCond) Match(node *Node) bool {
	if cond == nil {
		return true
	}
	for _, nmc := range cond {
		if nmc != nil && nmc.Match(node) {
			return true
		}
	}
	return false
}

// LinkMatchClause is a conjunction of conditions to match a semantic link.
//
// A semantic link satisfies the LinkMatchClause
// if it satisfies all the conditions in this LinkMatchClause.
//
// LinkMatchClause can specify the link ID, link type, properties on the link,
// the node from which the link starts, and the node to which the link points.
type LinkMatchClause interface {
	NLMatchClause

	// GetFromNodeMatchClause returns the match conditions for
	// the node from which the link starts.
	//
	// If there is no limit on the node, it returns nil.
	GetFromNodeMatchClause() NodeMatchClause

	// SetFromNodeMatchClause specifies the match conditions for
	// the node from which the link starts.
	//
	// If nmc is nil, it considers no limit on the node.
	SetFromNodeMatchClause(nmc NodeMatchClause)

	// GetToNodeMatchClause returns the match conditions for
	// the node to which the link points.
	//
	// If there is no limit on the node, it returns nil.
	GetToNodeMatchClause() NodeMatchClause

	// SetToNodeMatchClause specifies the match conditions for
	// the node to which the link points.
	//
	// If nmc is nil, it considers no limit on the node.
	SetToNodeMatchClause(nmc NodeMatchClause)

	// Match reports whether the semantic link satisfies this LinkMatchClause.
	Match(link *Link) bool
}

type linkMatchClauseImpl struct {
	nlMatchClauseImpl
	from NodeMatchClause // Match conditions for the node from which the link starts.
	to   NodeMatchClause // Match conditions for the node to which the link points.
}

// NewLinkMatchClause creates a new LinkMatchClause.
func NewLinkMatchClause() LinkMatchClause {
	return new(linkMatchClauseImpl)
}

func (lmc *linkMatchClauseImpl) SetIDAndClearOtherConds(id ID) {
	lmc.SetID(id)
	lmc.t, lmc.pmc, lmc.from, lmc.to = Type{}, nil, nil, nil
}

func (lmc *linkMatchClauseImpl) GetFromNodeMatchClause() NodeMatchClause {
	return lmc.from
}

func (lmc *linkMatchClauseImpl) SetFromNodeMatchClause(nmc NodeMatchClause) {
	lmc.from = nmc
}

func (lmc *linkMatchClauseImpl) GetToNodeMatchClause() NodeMatchClause {
	return lmc.to
}

func (lmc *linkMatchClauseImpl) SetToNodeMatchClause(nmc NodeMatchClause) {
	lmc.to = nmc
}

func (lmc *linkMatchClauseImpl) Match(link *Link) bool {
	switch {
	case link == nil:
	case lmc.id.IsValid() && link.ID != lmc.id:
	case lmc.t.IsValid() && link.Type != lmc.t:
	case lmc.pmc != nil && !lmc.pmc.Match(link.Props):
	case lmc.from != nil && !lmc.from.Match(link.From):
	case lmc.to != nil && !lmc.to.Match(link.To):
	default:
		return true
	}
	return false
}

// LinkMatchCond is a disjunction of the clauses of type LinkMatchClause
// to match a semantic link.
//
// Any nil LinkMatchClause in the LinkMatchCond is ignored.
//
// A semantic link satisfies the LinkMatchCond
// if it satisfies any of these clauses.
//
// In particular, a nil LinkMatchCond matches any semantic link (including nil).
// A non-nil but empty LinkMatchCond matches nothing.
type LinkMatchCond []LinkMatchClause

// Match reports whether the semantic link satisfies this LinkMatchCond.
func (cond LinkMatchCond) Match(link *Link) bool {
	if cond == nil {
		return true
	}
	for _, lmc := range cond {
		if lmc != nil && lmc.Match(link) {
			return true
		}
	}
	return false
}
