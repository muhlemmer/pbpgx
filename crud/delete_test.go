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

func TestTable_DeleteOne(t *testing.T) {
	tests := []struct {
		name    string
		query   *support.SimpleQuery
		want    *support.Simple
		wantErr bool
	}{
		{
			"No return",
			&support.SimpleQuery{
				Id: 5,
			},
			nil,
			false,
		},
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
			false,
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
			false,
		},
		{
			"no row error",
			&support.SimpleQuery{
				Id: 99,
				Columns: []support.SimpleColumns{
					support.SimpleColumns_id,
				},
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

			got, err := simpleRoTab.DeleteOne(testlib.CTX, tx, tt.query.Id, tt.query.Columns...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Table.DeleteOne() err = %v wantErr: %v", err, tt.wantErr)
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("Table.DeleteOne() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
