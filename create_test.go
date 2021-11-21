// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
)

func Test_createOneQuery(t *testing.T) {
	type args struct {
		schema string
		table  string
		req    proto.Message
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantN     int
		wantArgs  []interface{}
		wantErr   bool
	}{
		{
			"default schema",
			args{
				"",
				"simple",
				&support.Simple{
					Id:    100,
					Title: "foobar",
				},
			},
			`INSERT INTO "public"."simple" ("id", "title") VALUES ($1, $2) RETURNING "id";`,
			77,
			[]interface{}{
				&pgtype.Int4{
					Int:    100,
					Status: pgtype.Present,
				},
				&pgtype.Text{
					String: "foobar",
					Status: pgtype.Present,
				},
			},
			false,
		},
		{
			"repeated cardinality",
			args{
				"public",
				"simple",
				&support.Unsupported{
					Bl: []bool{true, false},
				},
			},
			"",
			0,
			nil,
			true,
		},
		{
			"unsupported type",
			args{
				"public",
				"simple",
				&support.Unsupported{
					Sup: &support.Supported{},
				},
			},
			"",
			0,
			nil,
			true,
		},
		{
			"all supported type",
			args{
				"custom",
				"support",
				&support.Supported{
					Bl:  true,
					I32: 1,
					I64: 2,
					F:   1.1,
					D:   2.2,
					S:   "foobar",
					Bt:  []byte("Hello, world!"),
					U32: 3,
					U64: 4,
				},
			},
			`INSERT INTO "custom"."support" ("bl", "i32", "i64", "f", "d", "s", "bt", "u32", "u64") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "id";`,
			146,
			[]interface{}{
				&pgtype.Bool{Bool: true, Status: pgtype.Present},
				&pgtype.Int4{Int: 1, Status: pgtype.Present},
				&pgtype.Int8{Int: 2, Status: pgtype.Present},
				&pgtype.Float4{Float: 1.1, Status: pgtype.Present},
				&pgtype.Float8{Float: 2.2, Status: pgtype.Present},
				&pgtype.Text{String: "foobar", Status: pgtype.Present},
				&pgtype.Bytea{Bytes: []byte("Hello, world!"), Status: pgtype.Present},
				&pgtype.Int4{Int: 3, Status: pgtype.Present},
				&pgtype.Int8{Int: 4, Status: pgtype.Present},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()

				gotQuery, gotN, gotArgs := createOneQuery(tt.args.schema, tt.args.table, tt.args.req)
				if gotQuery != tt.wantQuery {
					t.Errorf("createOneQuery() gotQuery =\n%v\nwant\n%v", gotQuery, tt.wantQuery)
				}
				if gotN != tt.wantN {
					t.Errorf("createOneQuery() gotN = %v, want %v", gotN, tt.wantN)
				}
				if len(gotArgs) != len(tt.wantArgs) {
					t.Errorf("createOneQuery() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
				}

				for i, want := range tt.wantArgs {
					if !reflect.DeepEqual(gotArgs[i], want) {
						t.Errorf("createOneQuery() gotArgs[%d] = %v, want %v", i, gotArgs[i], want)
					}
				}

				return nil
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("createOneQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Benchmark_createOneQuery(b *testing.B) {
	req := &support.Simple{
		Id:    100,
		Title: "foobar",
	}

	for i := 0; i < b.N; i++ {
		createOneQuery("public", "simple_rw", req)
	}
}

func TestCreateOne(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		req     proto.Message
		want    interface{}
		wantErr bool
	}{
		{
			"unsupported error",
			testCtx,
			&support.Unsupported{
				Bl: []bool{true, false},
			},
			nil,
			true,
		},
		{
			"succes",
			testCtx,
			&support.Simple{
				Id:    45,
				Title: "testing",
			},
			int32(45),
			false,
		},
		{
			"context error",
			errCtx,
			&support.Simple{
				Id:    45,
				Title: "testing",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateOne[*support.Simple](testCtx, connPool, "", "simple_rw", tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateOne() = %T(%v)., want %T(%v)", got, got, tt.want, tt.want)
			}
		})
	}
}
