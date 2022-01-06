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

package pbpgx

import (
	"context"
	"testing"

	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"google.golang.org/protobuf/proto"
)

func TestQuery(t *testing.T) {
	type args struct {
		ctx  context.Context
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []*support.Simple
		wantErr bool
	}{
		{
			"CTX error",
			args{
				testlib.ECTX,
				"select * from simple_ro;",
				nil,
			},
			nil,
			true,
		},
		{
			"Simple query",
			args{
				testlib.CTX,
				"select * from simple_ro;",
				nil,
			},
			[]*support.Simple{
				{Id: 1, Title: "one", Data: "foo bar"},
				{Id: 2, Title: "two"},
				{Id: 3, Data: "golden triangle"},
				{Id: 4, Title: "four", Data: "hello world"},
				{Id: 5, Title: "five", Data: "five is a four letter word"},
			},
			false,
		},
		{
			"Query with arguments",
			args{
				testlib.CTX,
				"select title from simple_ro where id in ($1, $2, $3);",
				[]interface{}{2, 3, 4},
			},
			[]*support.Simple{
				{Title: "two"},
				{},
				{Title: "four"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query[*support.Simple](tt.args.ctx, testlib.ConnPool, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}

			for i, want := range tt.want {
				if !proto.Equal(got[i], want) {
					t.Errorf("Query().[%d] = %v, want %v", i, got[i], want)
				}
			}
		})
	}
}

func TestQueryRow(t *testing.T) {
	type args struct {
		ctx  context.Context
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *support.Simple
		wantErr bool
	}{
		{
			"CTX error",
			args{
				testlib.ECTX,
				"select * from simple_ro;",
				nil,
			},
			nil,
			true,
		},
		{
			"No rows",
			args{
				testlib.CTX,
				"select title from simple_ro where id = $1;",
				[]interface{}{99},
			},
			nil,
			true,
		},
		{
			"Simple query",
			args{
				testlib.CTX,
				"select * from simple_ro;",
				nil,
			},
			&support.Simple{Id: 1, Title: "one", Data: "foo bar"},
			false,
		},
		{
			"Query with argument",
			args{
				testlib.CTX,
				"select title from simple_ro where id = $1;",
				[]interface{}{2},
			},
			&support.Simple{Title: "two"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryRow[*support.Simple](tt.args.ctx, testlib.ConnPool, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Unsupported error", func(t *testing.T) {
		_, err := QueryRow[*support.Unsupported](testlib.CTX, testlib.ConnPool, "select true as bool;")
		if err == nil {
			t.Error("Query() error is nil, want error")
		}
	})
}

func TestQueryStream(t *testing.T) {
	type args struct {
		ctx    context.Context
		stream *testServerStream[*support.Simple]
		sql    string
		args   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []*support.Simple
		wantErr bool
	}{
		{
			"CTX error",
			args{
				testlib.ECTX,
				&testServerStream[*support.Simple]{
					ctx: testlib.ECTX,
				},
				"select * from simple_ro;",
				nil,
			},
			nil,
			true,
		},
		{
			"Simple query",
			args{
				testlib.CTX,
				&testServerStream[*support.Simple]{
					ctx: testlib.CTX,
				},
				"select * from simple_ro;",
				nil,
			},
			[]*support.Simple{
				{Id: 1, Title: "one", Data: "foo bar"},
				{Id: 2, Title: "two"},
				{Id: 3, Data: "golden triangle"},
				{Id: 4, Title: "four", Data: "hello world"},
				{Id: 5, Title: "five", Data: "five is a four letter word"},
			},
			false,
		},
		{
			"Query with arguments",
			args{
				testlib.CTX,
				&testServerStream[*support.Simple]{
					ctx: testlib.CTX,
				},
				"select title from simple_ro where id in ($1, $2, $3);",
				[]interface{}{2, 3, 4},
			},
			[]*support.Simple{
				{Title: "two"},
				{},
				{Title: "four"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := QueryStream[*support.Simple](testlib.ConnPool, tt.args.stream, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := tt.args.stream.results

			if len(got) != len(tt.want) {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}

			for i, want := range tt.want {
				if !proto.Equal(got[i], want) {
					t.Errorf("Query().[%d] = %v, want %v", i, got[i], want)
				}
			}
		})
	}
}
