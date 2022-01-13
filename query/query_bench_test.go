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

func BenchmarkBuilder_Insert(b *testing.B) {
	type args struct {
		schema, table string
		insertColumns []string
		returnColumns []ColName
	}

	a := args{
		schema:        "public",
		table:         "simple",
		insertColumns: []string{"title", "data"},
		returnColumns: []ColName{
			support.SimpleColumns_id,
			support.SimpleColumns_title,
			support.SimpleColumns_data,
		},
	}

	var bldr Builder[ColName]

	for i := 0; i < b.N; i++ {
		bldr.Grow(bldr.lastCap)
		bldr.Insert(a.schema, a.table, a.insertColumns, a.returnColumns...)
		bldr.Reset()
	}
}

func BenchmarkBuilder_Select(b *testing.B) {
	type args struct {
		schema, table string
		columns       []support.SimpleColumns
		wf            WhereFunc[support.SimpleColumns]
		orderBy       OrderWriter[support.SimpleColumns]
		limit         int64
	}

	a := args{
		schema: "public",
		table:  "simple",
		columns: []support.SimpleColumns{
			support.SimpleColumns_title,
			support.SimpleColumns_data,
		},
		wf: WhereID[support.SimpleColumns],
		orderBy: Order(
			Ascending,
			support.SimpleColumns_title,
			support.SimpleColumns_id,
		),
		limit: 1,
	}

	var bldr Builder[support.SimpleColumns]

	for i := 0; i < b.N; i++ {
		bldr.Grow(bldr.lastCap)
		bldr.Select(a.schema, a.table, a.columns, a.wf, a.orderBy, a.limit)
		bldr.Reset()
	}
}

func BenchmarkBuilder_Update(b *testing.B) {
	type args struct {
		schema, table string
		wf            WhereFunc[ColName]
		updateColumns []string
		returnColumns []ColName
	}

	a := args{
		schema:        "public",
		table:         "simple",
		updateColumns: []string{"id", "title"},
		wf:            WhereID[ColName],
		returnColumns: []ColName{
			support.SimpleColumns_id,
			support.SimpleColumns_title,
			support.SimpleColumns_data,
		},
	}

	var bldr Builder[ColName]

	for i := 0; i < b.N; i++ {
		bldr.Grow(bldr.lastCap)
		bldr.Insert(a.schema, a.table, a.updateColumns, a.returnColumns...)
		bldr.Reset()
	}
}

func BenchmarkBuilder_Delete(b *testing.B) {
	type args struct {
		schema, table string
		wf            WhereFunc[ColName]
		returnColumns []ColName
	}

	a := args{
		schema: "public",
		table:  "simple",
		wf:     WhereID[ColName],
		returnColumns: []ColName{
			support.SimpleColumns_title,
			support.SimpleColumns_data,
		},
	}

	var bldr Builder[ColName]

	for i := 0; i < b.N; i++ {
		bldr.Grow(bldr.lastCap)
		bldr.Delete(a.schema, a.table, a.wf, a.returnColumns...)
		bldr.Reset()
	}
}
