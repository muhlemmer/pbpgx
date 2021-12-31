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

	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"google.golang.org/protobuf/proto"
)

func TestTable_ReadOne(t *testing.T) {
	tests := []struct {
		name    string
		id      int32
		want    *support.Simple
		wantErr bool
	}{
		{
			"not found",
			99,
			nil,
			true,
		},
		{
			"success",
			5,
			&support.Simple{
				Id:    5,
				Title: "five",
				Data:  "five is a four letter word",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := simpleRoTab.ReadOne(testlib.CTX, testlib.ConnPool, tt.id, []support.SimpleColumns{
				support.SimpleColumns_id,
				support.SimpleColumns_data,
				support.SimpleColumns_title,
			})

			if (err != nil) != tt.wantErr {
				t.Fatalf("Table.ReadOne() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("Table.ReadOne() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestTable_ReadAll(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    []*support.Simple
		wantErr bool
	}{
		{
			"context error",
			testlib.ECTX,
			nil,
			true,
		},
		{
			"success",
			testlib.CTX,
			[]*support.Simple{
				{
					Title: "one",
					Data:  "foo bar",
				},
				{
					Title: "two",
				},
				{
					Data: "golden triangle",
				},
				{
					Title: "four",
					Data:  "hello world",
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := simpleRoTab.ReadAll(tt.ctx, testlib.ConnPool, 4, []support.SimpleColumns{
				support.SimpleColumns_title,
				support.SimpleColumns_data,
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("Table.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("Table.ReadAll() =\n%v\nwant\n%v", got, tt.want)
			}

			for i, want := range tt.want {
				if !proto.Equal(got[i], want) {
					t.Fatalf("Table.ReadAll() =\n%v\nwant\n%v", got[i], want)
				}
			}
		})
	}
}

func TestTable_ReadList(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    []*support.Simple
		wantErr bool
	}{
		{
			"context error",
			testlib.ECTX,
			nil,
			true,
		},
		{
			"success",
			testlib.CTX,
			[]*support.Simple{
				{
					Title: "one",
					Data:  "foo bar",
				},
				{
					Title: "four",
					Data:  "hello world",
				},
				{
					Title: "five",
					Data:  "five is a four letter word",
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := simpleRoTab.ReadList(tt.ctx, testlib.ConnPool, []int32{1, 4, 5}, []support.SimpleColumns{
				support.SimpleColumns_title,
				support.SimpleColumns_data,
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("Table.ReadList() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("Table.ReadList() =\n%v\nwant\n%v", got, tt.want)
			}

			for i, want := range tt.want {
				if !proto.Equal(got[i], want) {
					t.Fatalf("Table.ReadList() =\n%v\nwant\n%v", got[i], want)
				}
			}
		})
	}
}
