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

import "testing"

func Test_whereID(t *testing.T) {
	b := &Builder[ColName]{
		argPos: 2,
	}

	WhereID(b)

	const want = " WHERE \"id\" = $3"

	if got := b.String(); got != want {
		t.Errorf("WhereID = %s, want %s", got, want)
	}
}

func Test_whereIDInFunc(t *testing.T) {
	b := &Builder[ColName]{
		argPos: 2,
	}

	WhereIDInFunc[ColName](3)(b)

	const want = " WHERE \"id\" IN ($3, $4, $5)"

	if got := b.String(); got != want {
		t.Errorf("whereID = %s, want %s", got, want)
	}
}
