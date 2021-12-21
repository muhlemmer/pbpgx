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

// WhereFunc are callback functions that writes
// the "WHERE" clause to a query.
type WhereFunc[Col ColName] func(b *Builder[Col])

// Where ID writes a where clause in the form:
//   WHERE "id" = $1
func WhereID[Col ColName](b *Builder[Col]) {
	b.WriteString(" WHERE \"id\" = ")
	b.WritePosArgs(1)
}

// WhereIDInFunc returns a function which writes a where clause in the form:
//   WHERE "id" IN $1, $2, $N...
func WhereIDInFunc[Col ColName](n int) WhereFunc[Col] {
	return func(b *Builder[Col]) {
		b.WriteString(" WHERE \"id\" IN ")
		b.WriteByte('(')
		b.WritePosArgs(n)
		b.WriteByte(')')
	}
}
