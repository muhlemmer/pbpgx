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
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"github.com/muhlemmer/pbpgx/query"
	"google.golang.org/protobuf/proto"
)

func TestCreateReturnOne(t *testing.T) {
	tests := []struct {
		name    string
		req     proto.Message
		retCols []support.SimpleColumns
		want    proto.Message
		wantErr bool
	}{
		{
			"all fields",
			&support.Simple{
				Id:    789,
				Title: "foo bar",
				Data:  "Hello World!",
			},
			[]support.SimpleColumns{
				support.SimpleColumns_id,
				support.SimpleColumns_title,
				support.SimpleColumns_data,
			},
			&support.Simple{
				Id:    789,
				Title: "foo bar",
				Data:  "Hello World!",
			},
			false,
		},
		{
			"unsupported error",
			&support.Unsupported{Bl: []bool{true, false}},
			[]support.SimpleColumns{
				support.SimpleColumns_id,
				support.SimpleColumns_title,
				support.SimpleColumns_data,
			},
			(*support.Simple)(nil),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(testlib.CTX, time.Second)
			defer cancel()

			tx, err := testlib.ConnPool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)

			got, err := simpleRwTab.CreateReturnOne(ctx, tx, tt.req, tt.retCols)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateReturnOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("CreateReturnOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateOne(t *testing.T) {
	tests := []struct {
		name    string
		req     proto.Message
		want    *support.Simple
		wantErr bool
	}{
		{
			"all fields",
			&support.Simple{
				Id:    991,
				Title: "foo bar",
				Data:  "Hello World!",
			},
			&support.Simple{
				Id:    991,
				Title: "foo bar",
				Data:  "Hello World!",
			},
			false,
		},
		{
			"unsupported error",
			&support.Unsupported{Bl: []bool{true, false}},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(testlib.CTX, time.Second)
			defer cancel()

			tx, err := testlib.ConnPool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)

			_, err = simpleRwTab.CreateOne(ctx, tx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateReturnOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := pbpgx.QueryRow[*support.Simple](ctx, tx, `select * from "simple_rw" where id = $1;`, tt.req.(*support.Simple).Id)
				if err != nil {
					t.Fatal(err)
				}

				if !proto.Equal(got, tt.want) {
					t.Errorf("CreateReturnOne() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func TestTable_Create(t *testing.T) {
	tests := []struct {
		name    string
		data    []proto.Message
		want    []pgconn.CommandTag
		wantErr bool
	}{
		{
			"no action",
			nil,
			nil,
			false,
		},
		{
			"all fields",
			[]proto.Message{
				&support.Simple{
					Id:    911,
					Title: "foo bar",
					Data:  "Hello World!",
				},
				&support.Simple{
					Id:    912,
					Title: "foo bar",
					Data:  "Hello World!",
				},
				&support.Simple{
					Id:    913,
					Title: "foo bar",
					Data:  "Hello World!",
				},
			},
			[]pgconn.CommandTag{
				[]byte("INSERT 0 1"),
				[]byte("INSERT 0 1"),
				[]byte("INSERT 0 1"),
			},
			false,
		},
		{
			"column mismatch error",
			[]proto.Message{
				&support.Supported{},
			},
			nil,
			true,
		},
		{
			"conflict error",
			[]proto.Message{
				&support.Simple{
					Id:    911,
					Title: "foo bar",
					Data:  "Hello World!",
				},
				&support.Simple{
					Id:    911,
					Title: "foo bar",
					Data:  "Hello World!",
				},
			},
			[]pgconn.CommandTag{
				[]byte("INSERT 0 1"),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(testlib.CTX, time.Second)
			defer cancel()

			tx, err := testlib.ConnPool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)

			got, err := simpleRwTab.Create(ctx, tx, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_Create_unsupported(t *testing.T) {
	tab := NewTable[query.ColName, *support.Unsupported, int32]("public", "unsupported", nil)

	ctx, cancel := context.WithTimeout(testlib.CTX, time.Second)
	defer cancel()

	tx, err := testlib.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	data := []proto.Message{
		&support.Unsupported{
			Bl: []bool{true},
		},
	}

	_, err = tab.Create(ctx, tx, data)
	if err == nil {
		t.Errorf("Create() error = %v, wantErr %v", err, true)
		return
	}
}
