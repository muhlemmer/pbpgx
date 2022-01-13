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

import (
	"testing"

	"github.com/muhlemmer/pbpgx/internal/support"
)

func TestOrder(t *testing.T) {
	tests := []struct {
		name      string
		direction Direction
		columns   []support.SimpleColumns
		want      string
	}{
		{
			"no columns",
			Ascending,
			nil,
			"",
		},
		{
			"one col, asc",
			Ascending,
			[]support.SimpleColumns{support.SimpleColumns_id},
			` ORDER BY "id" ASC`,
		},
		{
			"one col, desc",
			Descending,
			[]support.SimpleColumns{support.SimpleColumns_id},
			` ORDER BY "id" DESC`,
		},
		{
			"three col, asc",
			Ascending,
			[]support.SimpleColumns{
				support.SimpleColumns_id,
				support.SimpleColumns_title,
				support.SimpleColumns_created,
			},
			` ORDER BY "id", "title", "created" ASC`,
		},
		{
			"three col, desc",
			Descending,
			[]support.SimpleColumns{
				support.SimpleColumns_id,
				support.SimpleColumns_title,
				support.SimpleColumns_created,
			},
			` ORDER BY "id", "title", "created" DESC`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ow := Order(tt.direction, tt.columns...)
			var b Builder[support.SimpleColumns]
			ow.writeTo(&b)

			if got := b.String(); got != tt.want {
				t.Errorf("Order() = %s, want %s", got, tt.want)
			}
		})
	}
}
