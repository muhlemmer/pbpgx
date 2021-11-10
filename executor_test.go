// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
)

func TestExec(t *testing.T) {
	type args struct {
		ctx  context.Context
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    pgconn.CommandTag
		wantErr bool
	}{
		{
			"CTX error",
			args{
				errCtx,
				"insert into simple_rw (id, title) values ($1, $2);",
				[]interface{}{1, "foobar"},
			},
			nil,
			true,
		},
		{
			"Success",
			args{
				testCtx,
				"insert into simple_rw (id, title) values ($1, $2);",
				[]interface{}{1, "foobar"},
			},
			pgconn.CommandTag("INSERT 0 1"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exec(tt.args.ctx, connPool, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				errCtx,
				"select * from simple_ro;",
				nil,
			},
			nil,
			true,
		},
		{
			"Simple query",
			args{
				testCtx,
				"select * from simple_ro;",
				nil,
			},
			[]*support.Simple{
				{Id: 1, Title: "one"},
				{Id: 2, Title: "two"},
				{Id: 3},
				{Id: 4, Title: "four"},
				{Id: 5, Title: "five"},
			},
			false,
		},
		{
			"Query with arguments",
			args{
				testCtx,
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
			got, err := Query[*support.Simple](tt.args.ctx, connPool, tt.args.sql, tt.args.args...)
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
