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

func Test_scalarValue_PGValue(t *testing.T) {
	v := &scalarValue[bool]{
		ValueTranscoder: &pgtype.Bool{Status: pgtype.Present},
	}

	want := v.ValueTranscoder
	got := v.PGValue()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("listValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}

func Test_scalarValue_SetTo(t *testing.T) {
	got := &support.Supported{}
	msg := got.ProtoReflect()

	v := &scalarValue[int32]{
		ValueTranscoder: &pgtype.Int4{Status: pgtype.Present, Int: 2},
		fd:              msg.Descriptor().Fields().ByName("i32"),
		valueFunc:       pr.ValueOfInt32,
	}

	v.SetTo(msg)

	want := &support.Supported{I32: 2}

	if !proto.Equal(got, want) {
		t.Errorf("listValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}

func Test_scalarValue_SetTo_convert(t *testing.T) {
	got := &support.Supported{}
	msg := got.ProtoReflect()

	v := &scalarValue[int32]{
		ValueTranscoder: &pgtype.Int4{Status: pgtype.Present, Int: 2},
		fd:              msg.Descriptor().Fields().ByName("u32"),
		valueFunc:       convertIntValueFunc[int32](pr.ValueOfUint32),
	}

	v.SetTo(msg)

	want := &support.Supported{U32: 2}

	if !proto.Equal(got, want) {
		t.Errorf("listValue.PGValue =\n%v\nwant\n%v ", got, want)
	}
}
