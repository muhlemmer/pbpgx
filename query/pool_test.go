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
	"google.golang.org/protobuf/proto"
)

func TestPool(t *testing.T) {
	const want = `INSERT INTO "public"."simple" ("title", "data") VALUES ($1, $2) RETURNING "id", "title", "data";`

	type args struct {
		schema, table string
		msg           proto.Message
		returnColumns []ColName
		skipEmpty     bool
		exclude       []string
	}

	a := args{
		schema: "public",
		table:  "simple",
		msg: &support.Simple{
			Id:    31,
			Title: "foo bar",
		},
		returnColumns: []ColName{
			support.SimpleColumns_id,
			support.SimpleColumns_title,
			support.SimpleColumns_data,
		},
		skipEmpty: false,
		exclude:   []string{"id"},
	}

	var p Pool[ColName]

	for i := 0; i < 10; i++ {
		b := p.Get()
		b.Insert(a.schema, a.table, a.msg, a.returnColumns, a.skipEmpty, a.exclude...)

		if got := b.String(); got != want {
			t.Errorf("Builder.Insert() =\n%s\nwant\n%s", got, want)
		}

		p.Put(b)
	}
}