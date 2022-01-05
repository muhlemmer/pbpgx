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

func TestTable_CreateOne(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		req     proto.Message
		retCols []support.SimpleColumns
		want    proto.Message
		wantErr bool
	}{
		{
			"all fields, no return",
			testlib.CTX,
			&support.Simple{
				Id:    789,
				Title: "foo bar",
				Data:  "Hello World!",
			},
			nil,
			(*support.Simple)(nil),
			false,
		},
		{
			"all fields",
			testlib.CTX,
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
			testlib.CTX,
			&support.Unsupported{Bl: []bool{true, false}},
			[]support.SimpleColumns{
				support.SimpleColumns_id,
				support.SimpleColumns_title,
				support.SimpleColumns_data,
			},
			(*support.Simple)(nil),
			true,
		},
		{
			"context error",
			testlib.ECTX,
			&support.Simple{
				Id:    789,
				Title: "foo bar",
				Data:  "Hello World!",
			},
			nil,
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

			got, err := simpleRwTab.CreateOne(tt.ctx, tx, tt.req, tt.retCols)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table.CreateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("Table.CreateReturnOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_Create(t *testing.T) {
	tests := []struct {
		name    string
		cols    ColNames
		data    []proto.Message
		retCols []support.SimpleColumns
		want    []*support.Simple
		wantErr bool
	}{
		{
			"no action",
			nil,
			nil,
			nil,
			nil,
			false,
		},
		{
			"all fields, no return",
			ColNames{"id", "title", "data"},
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
			nil,
			nil,
			false,
		},
		{
			"all fields, return id",
			ColNames{"id", "title", "data"},
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
			[]support.SimpleColumns{
				support.SimpleColumns_id,
			},
			[]*support.Simple{
				{
					Id: 911,
				},
				{
					Id: 912,
				},
				{
					Id: 913,
				},
			},
			false,
		},
		{
			"column mismatch error",
			ColNames{"foo"},
			[]proto.Message{
				&support.Simple{
					Id:    911,
					Title: "foo bar",
					Data:  "Hello World!",
				},
			},
			nil,
			nil,
			true,
		},
		{
			"conflict error, return all",
			ColNames{"id", "title", "data"},
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
			[]support.SimpleColumns{
				support.SimpleColumns_id,
				support.SimpleColumns_data,
				support.SimpleColumns_title,
			},
			[]*support.Simple{
				{
					Id:    911,
					Title: "foo bar",
					Data:  "Hello World!",
				},
			},
			true,
		},
		{
			"unsupported error",
			ColNames{"bl"},
			[]proto.Message{
				&support.Unsupported{Bl: []bool{true, false}},
			},
			nil,
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

			got, err := simpleRwTab.Create(ctx, tx, tt.cols, tt.data, tt.retCols...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("Create() = %v, want %v", got, tt.want)
			}

			for i, want := range tt.want {
				if !proto.Equal(got[i], want) {
					t.Errorf("Create() = %v, want %v", got[i], want)
				}
			}
		})
	}
}

/*
func TestTable_Create_unsupported(t *testing.T) {
	tab := NewTable[support.SimpleColumns, *support.Unsupported, int32]("public", "unsupported", nil)

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
*/
