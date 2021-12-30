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
	"reflect"
	"testing"

	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
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
		msg           proto.Message
		returnColumns []ColName
		skipEmpty     bool
		exclude       []string
	}
	tests := []struct {
		name           string
		args           args
		want           string
		wantFieldNames []string
	}{
		{
			"no schema, skip empty, no return",
			args{
				schema: "",
				table:  "simple",
				msg: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				skipEmpty: true,
			},
			`INSERT INTO "simple" ("id", "title") VALUES ($1, $2);`,
			[]string{"id", "title"},
		},
		{
			"with schema, skip empty, return id",
			args{
				schema: "public",
				table:  "simple",
				msg: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				returnColumns: []ColName{
					support.SimpleColumns_id,
				},
				skipEmpty: true,
			},
			`INSERT INTO "public"."simple" ("id", "title") VALUES ($1, $2) RETURNING "id";`,
			[]string{"id", "title"},
		},
		{
			"with schema, no skip empty, exclude id, return all",
			args{
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
			},
			`INSERT INTO "public"."simple" ("title", "data") VALUES ($1, $2) RETURNING "id", "title", "data";`,
			[]string{"title", "data"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			fieldNames := b.Insert(tt.args.schema, tt.args.table, tt.args.msg, tt.args.returnColumns, tt.args.skipEmpty, tt.args.exclude...)

			if !reflect.DeepEqual(fieldNames, tt.wantFieldNames) {
				t.Errorf("Builder.Insert() fieldNames =\n%s\nwant\n%s", fieldNames, tt.wantFieldNames)
			}

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.Insert() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

func TestBuilder_Select(t *testing.T) {
	type args struct {
		schema, table string
		columns       []ColName
		wf            WhereFunc[ColName]
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
				columns: []ColName{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf:    WhereID[ColName],
				limit: 1,
			},
			`SELECT "title", "data" FROM "public"."simple" WHERE "id" = $1 LIMIT 1;`,
		},
		{
			"select all, with columns",
			args{
				schema: "public",
				table:  "simple",
				columns: []ColName{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf:    nil,
				limit: 0,
			},
			`SELECT "title", "data" FROM "public"."simple";`,
		},
		{
			"select list, with columns",
			args{
				schema: "public",
				table:  "simple",
				columns: []ColName{
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				wf:    WhereIDInFunc[ColName](3),
				limit: 0,
			},
			`SELECT "title", "data" FROM "public"."simple" WHERE "id" IN ($1, $2, $3);`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			b.Select(tt.args.schema, tt.args.table, tt.args.columns, tt.args.wf, tt.args.limit)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.Select() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Update(t *testing.T) {
	type args struct {
		schema, table string
		data          proto.Message
		wf            WhereFunc[ColName]
		returnColumns []ColName
		skipEmpty     bool
		exclude       []string
	}
	tests := []struct {
		name           string
		args           args
		want           string
		wantFieldNames []string
	}{
		{
			"no schema, skip empty, no return",
			args{
				schema: "",
				table:  "simple",
				data: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				wf:        WhereID[ColName],
				skipEmpty: true,
			},
			`UPDATE "simple" SET "id" = $1, "title" = $2 WHERE "id" = $3;`,
			[]string{"id", "title"},
		},
		{
			"with schema, skip empty, return all",
			args{
				schema: "public",
				table:  "simple",
				data: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				wf: WhereID[ColName],
				returnColumns: []ColName{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				skipEmpty: true,
			},
			`UPDATE "public"."simple" SET "id" = $1, "title" = $2 WHERE "id" = $3 RETURNING "id", "title", "data";`,
			[]string{"id", "title"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder[ColName]{}
			fieldNames := b.Update(tt.args.schema, tt.args.table, tt.args.data, tt.args.wf, tt.args.returnColumns, tt.args.skipEmpty)

			if !reflect.DeepEqual(fieldNames, tt.wantFieldNames) {
				t.Errorf("Builder.Update() fieldNames =\n%s\nwant\n%s", fieldNames, tt.wantFieldNames)
			}

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
			b.Delete(tt.args.schema, tt.args.table, tt.args.wf, tt.args.returnColumns)

			if got := b.String(); got != tt.want {
				t.Errorf("Builder.Delete() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}
