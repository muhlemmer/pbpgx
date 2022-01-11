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

type Value interface {
	pgtype.ValueTranscoder
	PGValue() pgtype.Value
	SetTo(pr.Message)
}

func New(fd pr.FieldDescriptor, status pgtype.Status) (v Value, err error) {
	repeated := fd.Cardinality() == pr.Repeated

	switch fd.Kind() {
	case pr.BoolKind:
		if repeated {
			v = &listValue[bool]{fd: fd, ValueTranscoder: &pgtype.BoolArray{Status: status}, valueFunc: pr.ValueOfBool}
		} else {
			v = &scalarValue[bool]{fd: fd, ValueTranscoder: &pgtype.Bool{Status: status}, valueFunc: pr.ValueOfBool}
		}

	case pr.Int32Kind, pr.Sint32Kind, pr.Sfixed32Kind:
		if repeated {
			v = &listValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4Array{Status: status}, valueFunc: pr.ValueOfInt32}
		} else {
			v = &scalarValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4{Status: status}, valueFunc: pr.ValueOfInt32}
		}

	case pr.Int64Kind, pr.Sint64Kind, pr.Sfixed64Kind:
		if repeated {
			v = &listValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8Array{Status: status}, valueFunc: pr.ValueOfInt64}
		} else {
			v = &scalarValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8{Status: status}, valueFunc: pr.ValueOfInt64}
		}

	case pr.FloatKind:
		if repeated {
			v = &listValue[float32]{fd: fd, ValueTranscoder: &pgtype.Float4Array{Status: status}, valueFunc: pr.ValueOfFloat32}
		} else {
			v = &scalarValue[float32]{fd: fd, ValueTranscoder: &pgtype.Float4{Status: status}, valueFunc: pr.ValueOfFloat32}
		}

	case pr.DoubleKind:
		if repeated {
			v = &listValue[float64]{fd: fd, ValueTranscoder: &pgtype.Float8Array{Status: status}, valueFunc: pr.ValueOfFloat64}
		} else {
			v = &scalarValue[float64]{fd: fd, ValueTranscoder: &pgtype.Float8{Status: status}, valueFunc: pr.ValueOfFloat64}
		}

	case pr.StringKind:
		if repeated {
			v = &listValue[string]{fd: fd, ValueTranscoder: &pgtype.TextArray{Status: status}, valueFunc: pr.ValueOfString}
		} else {
			v = &scalarValue[string]{fd: fd, ValueTranscoder: &pgtype.Text{Status: status}, valueFunc: pr.ValueOfString}
		}

	case pr.BytesKind:
		if repeated {
			v = &listValue[[]byte]{fd: fd, ValueTranscoder: &pgtype.ByteaArray{Status: status}, valueFunc: pr.ValueOfBytes}
		} else {
			v = &scalarValue[[]byte]{fd: fd, ValueTranscoder: &pgtype.Bytea{Status: status}, valueFunc: pr.ValueOfBytes}
		}

	case pr.Uint32Kind, pr.Fixed32Kind:
		if repeated {
			v = &listValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4Array{Status: status}, valueFunc: convertIntValueFunc[int32](pr.ValueOfUint32)}
		} else {
			v = &scalarValue[int32]{fd: fd, ValueTranscoder: &pgtype.Int4{Status: status}, valueFunc: convertIntValueFunc[int32](pr.ValueOfUint32)}
		}

	case pr.Uint64Kind, pr.Fixed64Kind:
		if repeated {
			v = &listValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8Array{Status: status}, valueFunc: convertIntValueFunc[int64](pr.ValueOfUint64)}
		} else {
			v = &scalarValue[int64]{fd: fd, ValueTranscoder: &pgtype.Int8{Status: status}, valueFunc: convertIntValueFunc[int64](pr.ValueOfUint64)}
		}

	case pr.MessageKind:
		name := fd.Message().FullName()

		switch name {
		case SupportedTimestamp:
			if repeated {
				v = &timestampListValue{fd: fd}
			} else {
				v = &timestampValue{fd: fd}
			}

		default:
			return nil, fmt.Errorf("unsupported message type %q", name)
		}

	default:
		return nil, fmt.Errorf("unsupported type %q", fd.Kind())
	}

	return v, nil
}
