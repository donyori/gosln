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

// PropMutateArg is the argument specifying
// how to mutate properties on a semantic node or link.
//
// It consists of two components:
//   - ToBeSet: a PropMap holding the properties to be set (added and replaced).
//   - ToBeRemoved: a PropNameSet holding the names of the properties to be removed.
//
// These two components are mutually exclusive:
// when a property is put into one component, it is removed from the other.
type PropMutateArg interface {
	// ToBeSet returns a PropMap with properties
	// to be set (added and replaced).
	//
	// The PropMap is always non-nil, but may be empty.
	ToBeSet() PropMap

	// ToBeRemoved returns a PropNameSet with
	// names of the properties to be removed.
	//
	// The PropNameSet is always non-nil, but may be empty.
	ToBeRemoved() PropNameSet
}

// propMutateArgImpl is an implementation of interface PropMutateArg.
type propMutateArgImpl struct {
	set    *mutExclPropMap     // Properties to set (add and replace).
	remove *mutExclPropNameSet // Names of the properties to remove.
}

// NewPropMutateArg creates a new PropMutateArg.
//
// setCap and removeCap ask to allocate enough space to hold
// the specified number of items in its ToBeSet and ToBeRemoved components,
// respectively.
// If setCap is negative, it is ignored, as is removeCap.
func NewPropMutateArg(setCap, removeCap int) PropMutateArg {
	pma := &propMutateArgImpl{
		set:    new(mutExclPropMap),
		remove: new(mutExclPropNameSet),
	}
	pma.set.init(setCap, pma.remove)
	pma.remove.init(removeCap, pma.set)
	return pma
}

func (pma *propMutateArgImpl) ToBeSet() PropMap {
	return pma.set
}

func (pma *propMutateArgImpl) ToBeRemoved() PropNameSet {
	return pma.remove
}
