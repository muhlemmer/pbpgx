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
	"testing"
	"time"

	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"google.golang.org/protobuf/proto"
)

func TestCreateReturnOne(t *testing.T) {
	tab := NewTable("", "simple_rw", ColumnDefault{"id": Zero})

	tests := []struct {
		name    string
		req     proto.Message
		retCols []support.SimpleColumns
		want    proto.Message
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateReturnOne[*support.Simple](testlib.CTX, testlib.ConnPool, tab, tt.req, tt.retCols)
			if err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("CreateReturnOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateOne(t *testing.T) {
	tab := NewTable("", "simple_rw", ColumnDefault{"id": Zero})

	tests := []struct {
		name string
		req  *support.Simple
		want *support.Simple
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

			_, err = CreateOne(testlib.CTX, testlib.ConnPool, tab, tt.req)
			if err != nil {
				t.Fatal(err)
			}

			got, err := pbpgx.QueryRow[*support.Simple](ctx, tx, `select * from "simple_rw" where id = $1;`, tt.req.Id)
			if err != nil {
				t.Fatal(err)
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("CreateReturnOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
