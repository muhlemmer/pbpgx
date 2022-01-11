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

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx/internal/support"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

func Test_listValue_PGValue(t *testing.T) {
	v := &listValue[bool]{
		ValueTranscoder: &pgtype.Bool{Status: pgtype.Present},
	}

	want := v.ValueTranscoder
	got := v.PGValue()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("listValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}

func Test_listValue_SetTo(t *testing.T) {
	got := &support.Supported{}
	msg := got.ProtoReflect()

	v := &listValue[int32]{
		ValueTranscoder: &pgtype.Int4Array{
			Status: pgtype.Present,
			Elements: []pgtype.Int4{
				{Status: pgtype.Present, Int: 2},
				{Status: pgtype.Present, Int: 4},
				{Status: pgtype.Present, Int: 6},
				{Status: pgtype.Present, Int: 8},
			},
		},
		fd:        msg.Descriptor().Fields().ByName("r_i32"),
		valueFunc: pr.ValueOfInt32,
	}

	v.SetTo(msg)

	want := &support.Supported{RI32: []int32{2, 4, 6, 8}}

	if !proto.Equal(got, want) {
		t.Errorf("listValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}
