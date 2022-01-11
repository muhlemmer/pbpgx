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
	"constraints"
	"fmt"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Value interface {
	pgtype.ValueTranscoder
	PGValue() pgtype.Value
	setTo(pr.Message)
}

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

func (v *scalarValue[T]) setTo(msg pr.Message) {
	if pgv, ok := v.Get().(T); ok {
		msg.Set(v.fd, v.valueFunc(pgv))
	}
}

type listValue[T scalar] struct {
	pgtype.ValueTranscoder
	fd        pr.FieldDescriptor
	valueFunc func(T) pr.Value
}

func (v *listValue[T]) PGValue() pgtype.Value { return v.ValueTranscoder }

func (v *listValue[T]) setTo(msg pr.Message) {
	var list []T
	v.AssignTo(&list)

	pl := msg.NewField(v.fd).List()
	for _, elem := range list {
		pl.Append(v.valueFunc(elem))
	}

	msg.Set(v.fd, pr.ValueOfList(pl))
}

// convertIntValueFunc returns a function which calls passed Value function
// with value of A converted the B.
func convertIntValueFunc[A, B constraints.Integer](f func(B) pr.Value) func(a A) pr.Value {
	return func(a A) pr.Value {
		return f(B(a))
	}
}

const (
	SupportedTimestamp = "google.protobuf.Timestamp"
)

type timestampValue struct {
	pgtype.Timestamptz
	fd pr.FieldDescriptor
}

func (v *timestampValue) PGValue() pgtype.Value { return &v.Timestamptz }

func (v *timestampValue) setTo(msg pr.Message) {
	msg.Set(v.fd, pr.ValueOfMessage(
		timestamppb.New(v.Time).ProtoReflect(),
	))
}

type timestampArrayValue struct {
	pgtype.TimestamptzArray
	fd pr.FieldDescriptor
}

func (v *timestampArrayValue) PGValue() pgtype.Value { return &v.TimestamptzArray }

func (v *timestampArrayValue) setTo(msg pr.Message) {
	pl := msg.NewField(v.fd).List()

	for _, x := range v.Elements {
		pl.Append(pr.ValueOfMessage(
			timestamppb.New(x.Time).ProtoReflect(),
		))
	}

	msg.Set(v.fd, pr.ValueOfList(pl))
}

func NewValue(fd pr.FieldDescriptor, status pgtype.Status) (v Value, err error) {
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
				v = &timestampArrayValue{fd: fd}
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
