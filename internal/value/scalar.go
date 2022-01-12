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
	"constraints"
	"fmt"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

type signedScalar interface {
	int32 | int64
}

type unsignedScalar interface {
	uint32 | uint64
}

type scalar interface {
	signedScalar | unsignedScalar | bool | float32 | float64 | string | []byte
}

type scalarValue[T scalar] struct {
	pgtype.ValueTranscoder
	fd        pr.FieldDescriptor
	valueFunc func(T) pr.Value
}

func (v *scalarValue[T]) PGValue() pgtype.Value { return v.ValueTranscoder }

func (v *scalarValue[T]) SetTo(msg pr.Message) {
	if pgv, ok := v.Get().(T); ok {
		msg.Set(v.fd, v.valueFunc(pgv))
	}
}

// convertIntValueFunc returns a function which calls passed Value function
// with value of A converted the B.
func convertIntValueFunc[A, B constraints.Integer](f func(B) pr.Value) func(a A) pr.Value {
	return func(a A) pr.Value {
		return f(B(a))
	}
}

func newScalarValue(fd pr.FieldDescriptor, status pgtype.Status) (v Value, err error) {
	switch fd.Kind() {
	case pr.BoolKind:
		v = &scalarValue[bool]{fd: fd, ValueTranscoder: &pgtype.Bool{Status: status}, valueFunc: pr.ValueOfBool}

	case pr.Int32Kind, pr.Sint32Kind, pr.Sfixed32Kind:
		v = &scalarValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4{Status: status}, valueFunc: pr.ValueOfInt32}

	case pr.Int64Kind, pr.Sint64Kind, pr.Sfixed64Kind:
		v = &scalarValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8{Status: status}, valueFunc: pr.ValueOfInt64}

	case pr.FloatKind:
		v = &scalarValue[float32]{fd: fd, ValueTranscoder: &pgtype.Float4{Status: status}, valueFunc: pr.ValueOfFloat32}

	case pr.DoubleKind:
		v = &scalarValue[float64]{fd: fd, ValueTranscoder: &pgtype.Float8{Status: status}, valueFunc: pr.ValueOfFloat64}

	case pr.StringKind:
		v = &scalarValue[string]{fd: fd, ValueTranscoder: &pgtype.Text{Status: status}, valueFunc: pr.ValueOfString}

	case pr.BytesKind:
		v = &scalarValue[[]byte]{fd: fd, ValueTranscoder: &pgtype.Bytea{Status: status}, valueFunc: pr.ValueOfBytes}

	case pr.Uint32Kind, pr.Fixed32Kind:
		v = &scalarValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4{Status: status}, valueFunc: convertIntValueFunc[int32](pr.ValueOfUint32)}

	case pr.Uint64Kind, pr.Fixed64Kind:
		v = &scalarValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8{Status: status}, valueFunc: convertIntValueFunc[int64](pr.ValueOfUint64)}

	default:
		return nil, fmt.Errorf("unsupported type %q", fd.Kind())
	}

	return v, nil
}
