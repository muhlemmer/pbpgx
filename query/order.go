/*
SPDX-License-Identifier: AGPL-3.0-only

Copyright (C) 2021, Tim Möhlmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package query

import "io"

// Direction used in a ORDER BY clause.
type Direction int

const (
	Ascending Direction = iota
	Descending
)

func (d Direction) writeTo(b io.StringWriter) {
	switch d {
	case Ascending:
		b.WriteString(" ASC")
	case Descending:
		b.WriteString(" DESC")
	}
}

type order[Col ColName] struct {
	Columns   []Col
	Direction Direction
}

func (o *order[Col]) writeTo(b *Builder[Col]) {
	if len(o.Columns) == 0 {
		return
	}

	const order = " ORDER BY "

	b.WriteString(order)
	b.WriteColumnSpec(o.Columns)
	o.Direction.writeTo(b)
}

// OrderWriter writen the ORDER BY ... [ASC|DESC] clause.
type OrderWriter[Col ColName] interface {
	writeTo(b *Builder[Col])
}

// Order returns an OrderWriter, which writes the
// ORDER BY "col1", "col2", "colN" [ASC|DESC] clause.
func Order[Col ColName](direction Direction, columns ...Col) OrderWriter[Col] {
	return &order[Col]{
		Columns:   columns,
		Direction: direction,
	}
}
