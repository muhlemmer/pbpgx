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

package value

import (
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_timestampValue_PGValue(t *testing.T) {
	v := &timestampValue{
		Timestamptz: pgtype.Timestamptz{
			Status: pgtype.Present,
			Time:   time.Unix(12, 0),
		},
	}

	want := &v.Timestamptz
	got := v.PGValue()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("timestampValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}

func Test_timestampValue_SetTo(t *testing.T) {
	got := &support.Supported{}
	msg := got.ProtoReflect()

	v := &timestampValue{
		Timestamptz: pgtype.Timestamptz{
			Status: pgtype.Present,
			Time:   time.Unix(12, 0),
		},
		fd: msg.Descriptor().Fields().ByName("ts"),
	}

	v.SetTo(msg)

	want := &support.Supported{Ts: &timestamppb.Timestamp{Seconds: 12}}

	if !proto.Equal(got, want) {
		t.Errorf("timestampValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}

func Test_timestampListValue_PGValue(t *testing.T) {
	v := &timestampListValue{
		TimestamptzArray: pgtype.TimestamptzArray{
			Status: pgtype.Present,
			Elements: []pgtype.Timestamptz{
				{
					Status: pgtype.Present,
					Time:   time.Unix(12, 0),
				},
				{
					Status: pgtype.Present,
					Time:   time.Unix(34, 0),
				},
				{
					Status: pgtype.Present,
					Time:   time.Unix(56, 0),
				},
			},
		},
	}

	want := &v.TimestamptzArray
	got := v.PGValue()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("timestampListValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}

func Test_timestampListValue_SetTo(t *testing.T) {
	got := &support.Supported{}
	msg := got.ProtoReflect()

	v := &timestampListValue{
		TimestamptzArray: pgtype.TimestamptzArray{
			Status: pgtype.Present,
			Elements: []pgtype.Timestamptz{
				{
					Status: pgtype.Present,
					Time:   time.Unix(12, 0),
				},
				{
					Status: pgtype.Present,
					Time:   time.Unix(34, 0),
				},
				{
					Status: pgtype.Present,
					Time:   time.Unix(56, 0),
				},
			},
		},
		fd: msg.Descriptor().Fields().ByName("r_ts"),
	}

	v.SetTo(msg)

	want := &support.Supported{RTs: []*timestamppb.Timestamp{
		{Seconds: 12},
		{Seconds: 34},
		{Seconds: 56},
	}}

	if !proto.Equal(got, want) {
		t.Errorf("timestampListValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}
