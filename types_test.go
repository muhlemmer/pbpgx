// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

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

func Test_pgKind_value(t *testing.T) {
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
			d := &pgKind{
				ValueTranscoder: tt.fields.ValueTranscoder,
				valueFunc:       tt.fields.valueFunc,
			}
			if got := d.value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgKind.value() = %v, want %v", got, tt.want)
			}
		})
	}
}
