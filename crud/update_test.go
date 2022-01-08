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

	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"google.golang.org/protobuf/proto"
)

func TestTable_UpdateOne(t *testing.T) {
	type args struct {
		query *support.SimpleQuery
		data  proto.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *support.Simple
		wantErr bool
	}{
		{
			"update title",
			args{
				&support.SimpleQuery{
					Id: 5,
					Columns: []support.SimpleColumns{
						support.SimpleColumns_id,
						support.SimpleColumns_title,
						support.SimpleColumns_data,
					},
				},
				&support.Simple{Title: "five-and-halve"},
			},
			&support.Simple{
				Id:    5,
				Title: "five-and-halve",
				Data:  "five is a four letter word",
			},
			false,
		},
		{
			"update title and data",
			args{
				&support.SimpleQuery{
					Id: 5,
					Columns: []support.SimpleColumns{
						support.SimpleColumns_id,
						support.SimpleColumns_title,
						support.SimpleColumns_data,
					},
				},
				&support.Simple{Title: "five-and-halve", Data: "history"},
			},
			&support.Simple{
				Id:    5,
				Title: "five-and-halve",
				Data:  "history",
			},
			false,
		},
		{
			"update title, return ID",
			args{
				&support.SimpleQuery{
					Id: 5,
					Columns: []support.SimpleColumns{
						support.SimpleColumns_id,
					},
				},
				&support.Simple{Title: "five-and-halve"},
			},
			&support.Simple{
				Id: 5,
			},
			false,
		},
		{
			"update title, no return",
			args{
				&support.SimpleQuery{
					Id: 5,
				},
				&support.Simple{Title: "five-and-halve"},
			},
			nil,
			false,
		},
		{
			"unsupported error",
			args{
				&support.SimpleQuery{
					Id: 5,
					Columns: []support.SimpleColumns{
						support.SimpleColumns_id,
					},
				},
				&support.Unsupported{Bl: []bool{true, false}},
			},
			nil,
			true,
		},
		{
			"update id, dup error",
			args{
				&support.SimpleQuery{
					Id: 5,
				},
				&support.Simple{Id: 4},
			},
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

			got, err := simpleRoTab.UpdateOne(testlib.CTX, tx, ParseFields(tt.args.data, true), tt.args.query.GetId(), tt.args.data, tt.args.query.GetColumns()...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table.UpdateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("Table.UpdateOne() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
