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
	"fmt"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

type listValue[T scalar] struct {
	pgtype.ValueTranscoder
	fd        pr.FieldDescriptor
	valueFunc func(T) pr.Value
}

func (v *listValue[T]) PGValue() pgtype.Value { return v.ValueTranscoder }

func (v *listValue[T]) SetTo(msg pr.Message) {
	var list []T
	v.AssignTo(&list)

	pl := msg.NewField(v.fd).List()
	for _, elem := range list {
		pl.Append(v.valueFunc(elem))
	}

	msg.Set(v.fd, pr.ValueOfList(pl))
}
func newlistValue(fd pr.FieldDescriptor, status pgtype.Status) (v Value, err error) {
	switch fd.Kind() {
	case pr.BoolKind:
		v = &listValue[bool]{fd: fd, ValueTranscoder: &pgtype.BoolArray{Status: status}, valueFunc: pr.ValueOfBool}

	case pr.Int32Kind, pr.Sint32Kind, pr.Sfixed32Kind:
		v = &listValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4Array{Status: status}, valueFunc: pr.ValueOfInt32}

	case pr.Int64Kind, pr.Sint64Kind, pr.Sfixed64Kind:
		v = &listValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8Array{Status: status}, valueFunc: pr.ValueOfInt64}

	case pr.FloatKind:
		v = &listValue[float32]{fd: fd, ValueTranscoder: &pgtype.Float4Array{Status: status}, valueFunc: pr.ValueOfFloat32}

	case pr.DoubleKind:
		v = &listValue[float64]{fd: fd, ValueTranscoder: &pgtype.Float8Array{Status: status}, valueFunc: pr.ValueOfFloat64}

	case pr.StringKind:
		v = &listValue[string]{fd: fd, ValueTranscoder: &pgtype.TextArray{Status: status}, valueFunc: pr.ValueOfString}

	case pr.BytesKind:
		v = &listValue[[]byte]{fd: fd, ValueTranscoder: &pgtype.ByteaArray{Status: status}, valueFunc: pr.ValueOfBytes}

	case pr.Uint32Kind, pr.Fixed32Kind:
		v = &listValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4Array{Status: status}, valueFunc: convertIntValueFunc[int32](pr.ValueOfUint32)}

	case pr.Uint64Kind, pr.Fixed64Kind:
		v = &listValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8Array{Status: status}, valueFunc: convertIntValueFunc[int64](pr.ValueOfUint64)}

	default:
		return nil, fmt.Errorf("unsupported type %q", fd.Kind())
	}

	return v, nil
}
