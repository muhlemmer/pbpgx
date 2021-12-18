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

	"github.com/jackc/pgx/v4"
	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"google.golang.org/protobuf/proto"
)

func TestDeleteOne(t *testing.T) {
	tab := NewTable("", "simple_ro", ColumnDefault{"id": Zero})

	tests := []struct {
		name string
		data proto.Message
	}{
		{
			"update title",
			&support.Simple{Title: "five-and-halve"},
		},
		{
			"update title and data",
			&support.Simple{Title: "five-and-halve", Data: "history"},
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

			_, err = DeleteOne(testlib.CTX, tx, tab, 5)
			if err != nil {
				t.Fatal(err)
			}

			_, err = pbpgx.QueryRow[*support.Simple](ctx, tx, `select * from "simple_ro" where id = $1;`, 5)
			if err != pgx.ErrNoRows {
				t.Errorf("DeleteOne() err = %v want %v", err, pgx.ErrNoRows)
			}
		})
	}
}

func TestDeleteReturnOne(t *testing.T) {
	tab := NewTable("", "simple_ro", ColumnDefault{"id": Zero})

	tests := []struct {
		name     string
		selector *support.SimpleQuery
		want     *support.Simple
	}{
		{
			"return all",
			&support.SimpleQuery{
				Id: 5,
				Columns: []support.SimpleColumns{
					support.SimpleColumns_id,
					support.SimpleColumns_title,
					support.SimpleColumns_data,
				},
			},
			&support.Simple{
				Id:    5,
				Title: "five",
				Data:  "five is a four letter word",
			},
		},
		{
			"return ID",
			&support.SimpleQuery{
				Id: 5,
				Columns: []support.SimpleColumns{
					support.SimpleColumns_id,
				},
			},
			&support.Simple{
				Id: 5,
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

			got, err := DeleteReturnOne[*support.Simple, int32, support.SimpleColumns](testlib.CTX, tx, tab, tt.selector)
			if err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("UpdateReturnOne() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
