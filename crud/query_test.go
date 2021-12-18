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

import (
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/stringx"
	"google.golang.org/protobuf/proto"
)

func Test_query_parseFields(t *testing.T) {
	req := &support.Simple{
		Id:    78,
		Title: "foo bar",
		Data:  "",
	}

	type args struct {
		skipEmpty bool
		exclude   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"no skip",
			args{false, nil},
			[]string{"id", "title", "data"},
		},
		{
			"exclude id",
			args{false, []string{"id"}},
			[]string{"title", "data"},
		},
		{
			"skip empty",
			args{true, nil},
			[]string{"id", "title"},
		},
		{
			"skip empty, exclude id",
			args{true, []string{"id"}},
			[]string{"title"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q query
			q.parseFields(req, tt.args.skipEmpty, tt.args.exclude...)

			if got := q.fieldNames; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.fieldNames =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func Test_query_parseArgs(t *testing.T) {
	type fields struct {
		fieldNames []string
		onEmpty    ColumnDefault
	}
	tests := []struct {
		name     string
		fields   fields
		req      proto.Message
		wantArgs []interface{}
	}{
		{
			"all empty, default",
			fields{
				fieldNames: []string{"id", "title", "data"},
			},
			&support.Simple{},
			[]interface{}{
				&pgtype.Int4{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Null},
			},
		},
		{
			"all empty, onEmpty map",
			fields{
				fieldNames: []string{"id", "title", "data"},
				onEmpty: ColumnDefault{
					"title": Zero,
					"data":  Null,
				},
			},
			&support.Simple{},
			[]interface{}{
				&pgtype.Int4{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Present},
				&pgtype.Text{Status: pgtype.Null},
			},
		},
		{
			"with id, onEmpty map",
			fields{
				fieldNames: []string{"id", "title", "data"},
				onEmpty: ColumnDefault{
					"title": Zero,
					"data":  Null,
				},
			},
			&support.Simple{
				Id: 87,
			},
			[]interface{}{
				&pgtype.Int4{Int: 87, Status: pgtype.Present},
				&pgtype.Text{Status: pgtype.Present},
				&pgtype.Text{Status: pgtype.Null},
			},
		},
		{
			"all values, onEmpty map",
			fields{
				fieldNames: []string{"id", "title", "data"},
				onEmpty: ColumnDefault{
					"title": Zero,
					"data":  Null,
				},
			},
			&support.Simple{
				Id:    87,
				Title: "Hello World!",
				Data:  "foo bar",
			},
			[]interface{}{
				&pgtype.Int4{Int: 87, Status: pgtype.Present},
				&pgtype.Text{String: "Hello World!", Status: pgtype.Present},
				&pgtype.Text{String: "foo bar", Status: pgtype.Present},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				fieldNames: tt.fields.fieldNames,
				tab: &Table{
					onEmpty: tt.fields.onEmpty,
				},
			}
			gotArgs := q.parseArgs(tt.req, 0)
			if len(gotArgs) != len(tt.wantArgs) {
				t.Fatalf("query.parseArgs() len = %v, want %v", gotArgs, tt.wantArgs)
			}

			for i := 0; i < len(gotArgs); i++ {
				got := gotArgs[i].(*pbpgx.Value).ValueTranscoder
				if !reflect.DeepEqual(got, tt.wantArgs[i]) {
					t.Errorf("query.parseArgs() = %v, want %v", gotArgs[i], tt.wantArgs[i])
				}
			}
		})
	}
}

func TestQuery_writeIdentifier(t *testing.T) {
	type fields struct {
		schema string
		table  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"w/o schema",
			fields{
				table: "bar",
			},
			"\"bar\"",
		},
		{
			"with schema",
			fields{
				schema: "foo",
				table:  "bar",
			},
			"\"foo\".\"bar\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewTable(tt.fields.schema, tt.fields.table, nil).newQuery(0)
			q.writeIdentifier()

			if got := q.String(); got != tt.want {
				t.Errorf("Table.writeIdentifier = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_query_writePosArgs(t *testing.T) {
	tests := []struct {
		name       string
		start, end int
		want       string
	}{
		{
			"one arg",
			1, 1,
			"$1",
		},
		{
			"no arg",
			1, 0,
			"",
		},
		{
			"multi arg",
			2, 4,
			"$2, $3, $4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				Builder: new(stringx.Builder),
			}

			q.writePosArgs(tt.start, tt.end)

			if got := q.String(); got != tt.want {
				t.Errorf("query.writePosArgs() = %s want %s", got, tt.want)
			}
		})
	}
}

func Test_insertQuery_noReturn(t *testing.T) {
	type args struct {
		tab       *Table
		req       proto.Message
		skipEmpty bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no schema, skip empty",
			args{
				tab: NewTable("", "simple", nil),
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				skipEmpty: true,
			},
			`INSERT INTO "simple" ("id", "title") VALUES ($1, $2);`,
		},
		{
			"with schema, skip empty",
			args{
				tab: NewTable("public", "simple", nil),
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				skipEmpty: true,
			},
			`INSERT INTO "public"."simple" ("id", "title") VALUES ($1, $2);`,
		},
		{
			"with schema, no skip empty",
			args{
				tab: NewTable("public", "simple", nil),
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				skipEmpty: false,
			},
			`INSERT INTO "public"."simple" ("id", "title", "data") VALUES ($1, $2, $3);`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := insertQuery(tt.args.tab, tt.args.req, 1000, NoRetCols, tt.args.skipEmpty)

			if got := q.String(); got != tt.want {
				t.Errorf("insertQuery() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

func Test_insertQuery_Return(t *testing.T) {
	type args struct {
		tab       *Table
		retCols   []support.SimpleColumns
		req       proto.Message
		skipEmpty bool
		exclude   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"exlude and return id",
			args{
				tab: NewTable("", "simple", nil),
				retCols: []support.SimpleColumns{
					support.SimpleColumns_id,
				},
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
					Data:  "Hello World!",
				},
				skipEmpty: false,
				exclude:   []string{"id"},
			},
			`INSERT INTO "simple" ("title", "data") VALUES ($1, $2) RETURNING "id";`,
		},
		{
			"exlude id and return all",
			args{
				tab: NewTable("", "simple", nil),
				retCols: []support.SimpleColumns{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
					Data:  "Hello World!",
				},
				skipEmpty: false,
				exclude:   []string{"id"},
			},
			`INSERT INTO "simple" ("title", "data") VALUES ($1, $2) RETURNING "id", "title", "data";`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := insertQuery(tt.args.tab, tt.args.req, 1000, tt.args.retCols, tt.args.skipEmpty, tt.args.exclude...)

			if got := q.String(); got != tt.want {
				t.Errorf("insertQuery() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

// INSERT INTO "public"."simple" ("id", "title", "data") VALUES ($1, $2, $3) RETURNING "id", "title", "data";
func Benchmark_insertQuery(b *testing.B) {
	tab := NewTable("public", "simple", nil)

	req := &support.Simple{
		Id:    31,
		Title: "foo bar",
		Data:  "Hello World!",
	}

	retCols := []support.SimpleColumns{
		support.SimpleColumns_id,
		support.SimpleColumns_title,
		support.SimpleColumns_data,
	}

	for i := 0; i < b.N; i++ {
		q := insertQuery(tab, req, tab.bufLens.createOne.get(), retCols, true)
		tab.bufLens.createOne.setHigher(q.Len())
		_ = q.String()
		q.release()
	}
}

func Test_selectQuery(t *testing.T) {
	tests := []struct {
		name string
		cols []support.SimpleColumns
		want string
	}{
		{
			"nil columns",
			nil,
			`SELECT * FROM "public"."simple_rw" WHERE "id" = $1 LIMIT 1;`,
		},
		{
			"with columns",
			[]support.SimpleColumns{
				support.SimpleColumns_title,
				support.SimpleColumns_data,
			},
			`SELECT "title", "data" FROM "public"."simple_rw" WHERE "id" = $1 LIMIT 1;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tab := NewTable("public", "simple_rw", nil)

			if got := selectQuery(tab, tt.cols, 1000).String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

// SELECT "id", "title", "data" FROM "public"."simple" WHERE "id" = $1 LIMIT 1;
func Benchmark_selectQuery(b *testing.B) {
	tab := NewTable("public", "simple", nil)

	retCols := []support.SimpleColumns{
		support.SimpleColumns_id,
		support.SimpleColumns_title,
		support.SimpleColumns_data,
	}

	for i := 0; i < b.N; i++ {
		q := selectQuery(tab, retCols, tab.bufLens.createOne.get())
		tab.bufLens.createOne.setHigher(q.Len())
		_ = q.String()
		q.release()
	}
}

func Test_updateQuery_noReturn(t *testing.T) {
	type args struct {
		tab       *Table
		req       proto.Message
		skipEmpty bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no schema, skip empty",
			args{
				tab: NewTable("", "simple", nil),
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				skipEmpty: true,
			},
			`UPDATE "simple" SET "id" = $1, "title" = $2 WHERE "id" = $3;`,
		},
		{
			"with schema, skip empty",
			args{
				tab: NewTable("public", "simple", nil),
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				skipEmpty: true,
			},
			`UPDATE "public"."simple" SET "id" = $1, "title" = $2 WHERE "id" = $3;`,
		},
		{
			"with schema, no skip empty",
			args{
				tab: NewTable("public", "simple", nil),
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
				},
				skipEmpty: false,
			},
			`UPDATE "public"."simple" SET "id" = $1, "title" = $2, "data" = $3 WHERE "id" = $4;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := updateQuery(tt.args.tab, tt.args.req, 1000, NoRetCols, tt.args.skipEmpty)

			if got := q.String(); got != tt.want {
				t.Errorf("updateQuery() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

func Test_updateQuery_Return(t *testing.T) {
	type args struct {
		tab       *Table
		retCols   []support.SimpleColumns
		req       proto.Message
		skipEmpty bool
		exclude   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"exlude and return id",
			args{
				tab: NewTable("", "simple", nil),
				retCols: []support.SimpleColumns{
					support.SimpleColumns_id,
				},
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
					Data:  "Hello World!",
				},
				skipEmpty: false,
				exclude:   []string{"id"},
			},
			`UPDATE "simple" SET "title" = $1, "data" = $2 WHERE "id" = $3 RETURNING "id";`,
		},
		{
			"exlude id and return all",
			args{
				tab: NewTable("", "simple", nil),
				retCols: []support.SimpleColumns{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
				req: &support.Simple{
					Id:    31,
					Title: "foo bar",
					Data:  "Hello World!",
				},
				skipEmpty: false,
				exclude:   []string{"id"},
			},
			`UPDATE "simple" SET "title" = $1, "data" = $2 WHERE "id" = $3 RETURNING "id", "title", "data";`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := updateQuery(tt.args.tab, tt.args.req, 1000, tt.args.retCols, tt.args.skipEmpty, tt.args.exclude...)

			if got := q.String(); got != tt.want {
				t.Errorf("updateQuery() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

// UPDATE "public"."simple" SET "id" = $1, "title" = $2, "data" = $3 WHERE "id" = $4 RETURNING "id", "title", "data";
func Benchmark_updateQuery(b *testing.B) {
	tab := NewTable("public", "simple", nil)

	data := &support.Simple{
		Id:    31,
		Title: "foo bar",
		Data:  "Hello World!",
	}

	retCols := []support.SimpleColumns{
		support.SimpleColumns_id,
		support.SimpleColumns_title,
		support.SimpleColumns_data,
	}

	for i := 0; i < b.N; i++ {
		q := updateQuery(tab, data, tab.bufLens.createOne.get(), retCols, true)
		tab.bufLens.createOne.setHigher(q.Len())
		_ = q.String()
		q.release()
	}
}

func Test_deleteQuery_noReturn(t *testing.T) {
	tests := []struct {
		name string
		tab  *Table
		want string
	}{
		{
			"no schema",
			NewTable("", "simple", nil),
			`DELETE FROM "simple" WHERE "id" = $1;`,
		},
		{
			"with schema",
			NewTable("public", "simple", nil),
			`DELETE FROM "public"."simple" WHERE "id" = $1;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := deleteQuery(tt.tab, 1000, NoRetCols)

			if got := q.String(); got != tt.want {
				t.Errorf("deleteQuery() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

func Test_deleteQuery_Return(t *testing.T) {
	type args struct {
		tab     *Table
		retCols []support.SimpleColumns
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"return id",
			args{
				tab: NewTable("", "simple", nil),
				retCols: []support.SimpleColumns{
					support.SimpleColumns_id,
				},
			},
			`DELETE FROM "simple" WHERE "id" = $1 RETURNING "id";`,
		},
		{
			"return all",
			args{
				tab: NewTable("", "simple", nil),
				retCols: []support.SimpleColumns{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
			},
			`DELETE FROM "simple" WHERE "id" = $1 RETURNING "id", "title", "data";`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := deleteQuery(tt.args.tab, 1000, tt.args.retCols)

			if got := q.String(); got != tt.want {
				t.Errorf("deleteQuery() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

// DELETE FROM "public"."simple" WHERE "id" = $1 RETURNING "id", "title", "data";
func Benchmark_deleteQuery(b *testing.B) {
	tab := NewTable("public", "simple", nil)

	retCols := []support.SimpleColumns{
		support.SimpleColumns_id,
		support.SimpleColumns_title,
		support.SimpleColumns_data,
	}

	for i := 0; i < b.N; i++ {
		q := deleteQuery(tab, tab.bufLens.createOne.get(), retCols)
		tab.bufLens.createOne.setHigher(q.Len())
		_ = q.String()
		q.release()
	}
}
