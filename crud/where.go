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

package crud

type whereFunc func(q *query)

func whereID(q *query) {
	q.WriteString(" WHERE \"id\" = ")
	q.writePosArgs(1)
}

func whereIDInFunc(n int) whereFunc {
	return func(q *query) {
		q.WriteString(" WHERE \"id\" IN ")
		q.WriteByte('(')
		q.writePosArgs(n)
		q.WriteByte(')')
	}
}
