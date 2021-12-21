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

package query

import (
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
)

func TestParseFields(t *testing.T) {
	msg := &support.Simple{
		Id:    78,
		Title: "foo bar",
		Data:  "",
	}

	type args struct {
		skipEmpty bool
		exclude   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"no skip",
			args{false, nil},
			[]string{"id", "title", "data"},
		},
		{
			"exclude id",
			args{false, []string{"id"}},
			[]string{"title", "data"},
		},
		{
			"skip empty",
			args{true, nil},
			[]string{"id", "title"},
		},
		{
			"skip empty, exclude id",
			args{true, []string{"id"}},
			[]string{"title"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseFields(msg, tt.args.skipEmpty, tt.args.exclude...)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.fieldNames =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	type args struct {
		msg        proto.Message
		fieldNames []string
		cd         ColumnDefault
	}
	tests := []struct {
		name     string
		args     args
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			"all empty, default",
			args{
				msg:        &support.Simple{},
				fieldNames: []string{"id", "title", "data"},
			},
			[]interface{}{
				&pgtype.Int4{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Null},
			},
			false,
		},
		{
			"all empty, onEmpty map",
			args{
				msg:        &support.Simple{},
				fieldNames: []string{"id", "title", "data"},
				cd: ColumnDefault{
					"title": Zero,
					"data":  Null,
				},
			},
			[]interface{}{
				&pgtype.Int4{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Present},
				&pgtype.Text{Status: pgtype.Null},
			},
			false,
		},
		{
			"with id, onEmpty map",
			args{
				msg: &support.Simple{
					Id: 87,
				},
				fieldNames: []string{"id", "title", "data"},
				cd: ColumnDefault{
					"title": Zero,
					"data":  Null,
				},
			},
			[]interface{}{
				&pgtype.Int4{Int: 87, Status: pgtype.Present},
				&pgtype.Text{Status: pgtype.Present},
				&pgtype.Text{Status: pgtype.Null},
			},
			false,
		},
		{
			"all values, onEmpty map",
			args{
				msg: &support.Simple{
					Id:    87,
					Title: "Hello World!",
					Data:  "foo bar",
				},
				fieldNames: []string{"id", "title", "data"},
				cd: ColumnDefault{
					"title": Zero,
					"data":  Null,
				},
			},
			[]interface{}{
				&pgtype.Int4{Int: 87, Status: pgtype.Present},
				&pgtype.Text{String: "Hello World!", Status: pgtype.Present},
				&pgtype.Text{String: "foo bar", Status: pgtype.Present},
			},
			false,
		},
		{
			"unsupported error",
			args{
				msg: &support.Unsupported{
					Bl: []bool{true, false},
				},
				fieldNames: []string{"bl"},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotArgs, err := ParseArgs(tt.args.msg, tt.args.fieldNames, tt.args.cd)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotArgs) != len(tt.wantArgs) {
				t.Fatalf("query.parseArgs() len = %v, want %v", gotArgs, tt.wantArgs)
			}

			for i := 0; i < len(gotArgs); i++ {
				got := gotArgs[i].(*pbpgx.Value).ValueTranscoder
				if !reflect.DeepEqual(got, tt.wantArgs[i]) {
					t.Errorf("query.parseArgs() = %v, want %v", gotArgs[i], tt.wantArgs[i])
				}
			}
		})
	}
}
