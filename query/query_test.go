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

func TestBuilder_Reset(t *testing.T) {
	b := &Builder[ColName]{}
	b.WritePosArgs(3)
	b.Reset()

	if b.argPos != 0 {
		t.Errorf("Builder.Reset argPos = %d, want %d", b.argPos, 0)
	}

	if b.Len() != 0 {
		t.Errorf("Builder.Reset Len = %d, want %d", b.Len(), 0)
	}
}

func TestBuilder_WriteIdentifier(t *testing.T) {
	type args struct {
		schema string
		table  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"w/o schema",
			args{
				table: "bar",
			},
			"\"bar\"",
		},
		{
			"with schema",
			args{
				schema: "foo",
				table:  "bar",
			},
			"\"foo\".\"bar\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			b.WriteIdentifier(tt.args.schema, tt.args.table)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.WriteIdentifier = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestBuilder_WritePosArgs(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{
			"one arg",
			1,
			"$1",
		},
		{
			"no arg",
			0,
			"",
		},
		{
			"multi arg",
			4,
			"$1, $2, $3, $4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			b.WritePosArgs(tt.n)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.WritePosArgs() = %s want %s", got, tt.want)
			}
		})
	}
}

func TestBuilder_WriteReturnClause(t *testing.T) {
	tests := []struct {
		name    string
		columns []ColName
		want    string
	}{
		{
			"no return clause",
			nil,
			``,
		},
		{
			"one column",
			[]ColName{
				support.SimpleColumns_id,
			},
			` RETURNING "id"`,
		},
		{
			"three columns",
			[]ColName{
				support.SimpleColumns_id,
				support.SimpleColumns_title,
				support.SimpleColumns_data,
			},
			` RETURNING "id", "title", "data"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			b.WriteReturnClause(tt.columns)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.WriteReturnClause() = %s want %s", got, tt.want)
			}
		})
	}
}

func TestBuilder_Insert(t *testing.T) {
	type args struct {
		schema, table string
		insertColumns []string
		returnColumns []ColName
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no schema, no return",
			args{
				schema:        "",
				table:         "simple",
				insertColumns: []string{"id", "title"},
			},
			`INSERT INTO "simple" ("id", "title") VALUES ($1, $2);`,
		},
		{
			"with schema, return id",
			args{
				schema:        "public",
				table:         "simple",
				insertColumns: []string{"id", "title"},
				returnColumns: []ColName{
					support.SimpleColumns_id,
				},
			},
			`INSERT INTO "public"."simple" ("id", "title") VALUES ($1, $2) RETURNING "id";`,
		},
		{
			"with schema, return all",
			args{
				schema:        "public",
				table:         "simple",
				insertColumns: []string{"title", "data"},
				returnColumns: []ColName{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
			},
			`INSERT INTO "public"."simple" ("title", "data") VALUES ($1, $2) RETURNING "id", "title", "data";`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			b.Insert(tt.args.schema, tt.args.table, tt.args.insertColumns, tt.args.returnColumns...)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.Insert() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

func TestBuilder_Select(t *testing.T) {
	type args struct {
		schema, table string
		columns       []support.SimpleColumns
		wf            WhereFunc[support.SimpleColumns]
		orderBy       OrderWriter[support.SimpleColumns]
		limit         int64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"select one, with columns",
			args{
				schema: "public",
				table:  "simple",
				columns: []support.SimpleColumns{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf:    WhereID[support.SimpleColumns],
				limit: 1,
			},
			`SELECT "title", "data" FROM "public"."simple" WHERE "id" = $1 LIMIT 1;`,
		},
		{
			"select all, with columns",
			args{
				schema: "public",
				table:  "simple",
				columns: []support.SimpleColumns{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf:    nil,
				limit: 0,
			},
			`SELECT "title", "data" FROM "public"."simple";`,
		},
		{
			"select all, with columns, ordered Ascending",
			args{
				schema: "public",
				table:  "simple",
				columns: []support.SimpleColumns{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf: nil,
				orderBy: Order(
					Ascending,
					support.SimpleColumns_title,
					support.SimpleColumns_id,
				),
				limit: 0,
			},
			`SELECT "title", "data" FROM "public"."simple" ORDER BY "title", "id" ASC;`,
		},
		{
			"select list, with columns",
			args{
				schema: "public",
				table:  "simple",
				columns: []support.SimpleColumns{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf:    WhereIDInFunc[support.SimpleColumns](3),
				limit: 0,
			},
			`SELECT "title", "data" FROM "public"."simple" WHERE "id" IN ($1, $2, $3);`,
		},
		{
			"select list, with columns, order asc",
			args{
				schema: "public",
				table:  "simple",
				columns: []support.SimpleColumns{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf: WhereIDInFunc[support.SimpleColumns](3),
				orderBy: Order(
					Ascending,
					support.SimpleColumns_id,
				),
				limit: 0,
			},
			`SELECT "title", "data" FROM "public"."simple" WHERE "id" IN ($1, $2, $3) ORDER BY "id" ASC;`,
		},
		{
			"select list, with columns, order desc",
			args{
				schema: "public",
				table:  "simple",
				columns: []support.SimpleColumns{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf: WhereIDInFunc[support.SimpleColumns](3),
				orderBy: Order(
					Descending,
					support.SimpleColumns_id,
				),
				limit: 0,
			},
			`SELECT "title", "data" FROM "public"."simple" WHERE "id" IN ($1, $2, $3) ORDER BY "id" DESC;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[support.SimpleColumns]{}
			b.Select(tt.args.schema, tt.args.table, tt.args.columns, tt.args.wf, tt.args.orderBy, tt.args.limit)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.Select() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Update(t *testing.T) {
	type args struct {
		schema, table string
		wf            WhereFunc[ColName]
		updateColumns []string
		returnColumns []ColName
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no schema, no return",
			args{
				schema:        "",
				table:         "simple",
				updateColumns: []string{"title"},
				wf:            WhereID[ColName],
			},
			`UPDATE "simple" SET "title" = $1 WHERE "id" = $2;`,
		},
		{
			"with schema, return all",
			args{
				schema:        "public",
				table:         "simple",
				updateColumns: []string{"id", "title"},
				wf:            WhereID[ColName],
				returnColumns: []ColName{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
			},
			`UPDATE "public"."simple" SET "id" = $1, "title" = $2 WHERE "id" = $3 RETURNING "id", "title", "data";`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			b.Update(tt.args.schema, tt.args.table, tt.args.updateColumns, tt.args.wf, tt.args.returnColumns...)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.Update() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

func TestBuilder_Delete(t *testing.T) {
	type args struct {
		schema, table string
		wf            WhereFunc[ColName]
		returnColumns []ColName
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no schema, no return",
			args{
				schema: "",
				table:  "simple",
				wf:     WhereID[ColName],
			},
			`DELETE FROM "simple" WHERE "id" = $1;`,
		},
		{
			"with schema, return all",
			args{
				schema: "public",
				table:  "simple",
				wf:     WhereID[ColName],
				returnColumns: []ColName{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
			},
			`DELETE FROM "public"."simple" WHERE "id" = $1 RETURNING "id", "title", "data";`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			b.Delete(tt.args.schema, tt.args.table, tt.args.wf, tt.args.returnColumns...)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.Delete() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}
