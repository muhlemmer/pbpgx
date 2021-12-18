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

package pbpgx

import (
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

func Test_convertIntValueFunc(t *testing.T) {
	f := convertIntValueFunc[uint32](pr.ValueOfInt32)
	var u uint32 = 222
	v := f(u)

	if got := v.Interface().(int32); got != 222 {
		t.Errorf("convertIntValueFunc.f() got: %v, want 222", got)
	}
}

func Test_Value_value(t *testing.T) {
	type fields struct {
		ValueTranscoder pgtype.ValueTranscoder
		valueFunc       func(interface{}) pr.Value
	}
	tests := []struct {
		name   string
		fields fields
		want   pr.Value
	}{
		{
			"test",
			fields{
				ValueTranscoder: &pgtype.Float4{
					Float:  1.1,
					Status: pgtype.Present,
				},
				valueFunc: pr.ValueOf,
			},
			pr.ValueOfFloat32(1.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Value{
				ValueTranscoder: tt.fields.ValueTranscoder,
				valueFunc:       tt.fields.valueFunc,
			}
			if got := d.value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgKind.value() = %v, want %v", got, tt.want)
			}
		})
	}
}
