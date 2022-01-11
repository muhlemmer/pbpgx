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
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/value"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseFields(t *testing.T) {
	msg := &support.Simple{
		Id:    78,
		Title: "foo bar",
		Data:  "",
	}

	tests := []struct {
		name      string
		skipEmpty bool
		ignore    []string
		want      ColNames
	}{
		{
			"no skip",
			false,
			nil,
			ColNames{"id", "title", "data", "created"},
		},
		{
			"skip empty",
			true,
			nil,
			ColNames{"id", "title"},
		},
		{
			"ignore id, title",
			false,
			[]string{"id", "title"},
			ColNames{"data", "created"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseFields(msg, tt.skipEmpty, tt.ignore...)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.fieldNames =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestColumns_ParseArgs(t *testing.T) {
	type args struct {
		msg  proto.Message
		cols ColNames
	}
	tests := []struct {
		name     string
		columns  Columns
		args     args
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			"all empty, default",
			nil,
			args{
				msg:  &support.Simple{},
				cols: []string{"id", "title", "data", "created"},
			},
			[]interface{}{
				&pgtype.Int4{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Null},
				&pgtype.Timestamptz{Status: pgtype.Null},
			},
			false,
		},
		{
			"all empty, onEmpty map",
			Columns{
				"title": Zero,
				"data":  Null,
			},
			args{
				msg:  &support.Simple{},
				cols: []string{"id", "title", "data", "created"},
			},
			[]interface{}{
				&pgtype.Int4{Status: pgtype.Null},
				&pgtype.Text{Status: pgtype.Present},
				&pgtype.Text{Status: pgtype.Null},
				&pgtype.Timestamptz{Status: pgtype.Null},
			},
			false,
		},
		{
			"with id, onEmpty map",
			Columns{
				"title": Zero,
				"data":  Null,
			},
			args{
				msg: &support.Simple{
					Id: 87,
				},
				cols: []string{"id", "title", "data"},
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
			Columns{
				"title": Zero,
				"data":  Null,
			},
			args{
				msg: &support.Simple{
					Id:    87,
					Title: "Hello World!",
					Data:  "foo bar",
					Created: &timestamppb.Timestamp{
						Seconds: 12,
						Nanos:   34,
					},
				},
				cols: []string{"id", "title", "data", "created"},
			},
			[]interface{}{
				&pgtype.Int4{Int: 87, Status: pgtype.Present},
				&pgtype.Text{String: "Hello World!", Status: pgtype.Present},
				&pgtype.Text{String: "foo bar", Status: pgtype.Present},
				&pgtype.Timestamptz{
					Time:   time.Unix(12, 34),
					Status: pgtype.Present,
				},
			},
			false,
		},
		{
			"unsupported error",
			nil,
			args{
				msg:  &support.Unsupported{Sup: &support.Supported{Bl: true}},
				cols: []string{"bl"},
			},
			nil,
			true,
		},
		{
			"unsupported message error",
			nil,
			args{
				msg: &support.Unsupported{
					Sup: &support.Supported{I32: 2},
				},
				cols: []string{"sup"},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotArgs, err := tt.columns.ParseArgs(tt.args.msg, tt.args.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotArgs) != len(tt.wantArgs) {
				t.Fatalf("query.parseArgs() len = %v, want %v", gotArgs, tt.wantArgs)
			}

			for i := 0; i < len(gotArgs); i++ {
				var ok bool
				got := gotArgs[i].(value.Value).PGValue()

				switch x := got.(type) {
				case *pgtype.Timestamptz:
					ok = x.Time.Equal(tt.wantArgs[i].(*pgtype.Timestamptz).Time)

				default:
					ok = reflect.DeepEqual(x, tt.wantArgs[i])
				}

				if !ok {
					t.Errorf("query.parseArgs() = %v, want %v", got, tt.wantArgs[i])
				}

			}
		})
	}
}
