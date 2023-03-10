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

package neo4jsln

import (
	"github.com/donyori/gogo/container/mapping"
	"github.com/donyori/gogo/errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/donyori/gosln"
)

// slnIDPropName is the property name of SLN ID in Cypher.
const slnIDPropName = "slnID"

// makeParameterMap renders a semantic node or link ID, a property map,
// and property names about to be removed as a parameter map for Cypher.
//
// If paraName is empty, makeParameterMap reports an error.
//
// If id is invalid, it is ignored.
func makeParameterMap(
	paraName string,
	id gosln.ID,
	props gosln.PropMap,
	remove gosln.PropNameSet,
) (para map[string]any, err error) {
	if paraName == "" {
		return nil, errors.AutoNew("parameter name is empty")
	}
	var n int
	if id.IsValid() {
		n = 1
	}
	if props != nil {
		n += props.Len()
	}
	if remove != nil {
		n += remove.Len()
	}
	if n == 0 {
		return map[string]any{paraName: nil}, nil
	}
	m := make(map[string]any, n)
	if id.IsValid() {
		m[slnIDPropName] = id.String()
	}
	if props != nil {
		props.Range(func(x mapping.Entry[gosln.PropName, any]) (cont bool) {
			if date, ok := x.Value.(gosln.Date); ok {
				m[x.Key.String()] = neo4j.DateOf(date.GoTime())
			} else {
				m[x.Key.String()] = x.Value
			}
			return true
		})
	}
	if remove != nil {
		remove.Range(func(x gosln.PropName) (cont bool) {
			m[x.String()] = nil
			return true
		})
	}
	return map[string]any{paraName: m}, nil
}
