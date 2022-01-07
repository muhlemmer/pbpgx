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
	"time"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertIntValueFunc returns a function which calls passed function
// with src asserted to the type of A and then converted to the type of B.
func convertIntValueFunc[A, B constraints.Integer](f func(B) pr.Value) func(interface{}) pr.Value {
	return func(src interface{}) pr.Value {
		v, _ := src.(A)
		return f(B(v))
	}
}

func assertValueFunc[T any](f func(T) pr.Value) func(interface{}) pr.Value {
	return func(src interface{}) pr.Value {
		v, _ := src.(T)
		return f(v)
	}
}

type Value struct {
	pgtype.ValueTranscoder
	fieldDesc pr.FieldDescriptor
	valueFunc func(interface{}) pr.Value
}

func (d *Value) value() pr.Value {
	return d.valueFunc(d.Get())
}

func NewValue(fd pr.FieldDescriptor, status pgtype.Status) (*Value, error) {
	d := &Value{
		fieldDesc: fd,
	}

	if fd.Cardinality() == pr.Repeated {
		return nil, fmt.Errorf("unsupported type \"%s %s\" for scanning", pr.Repeated, fd.Kind())
	}

	switch fd.Kind() {
	case pr.BoolKind:
		d.ValueTranscoder = &pgtype.Bool{Status: status}
		d.valueFunc = assertValueFunc(pr.ValueOfBool)

	case pr.Int32Kind, pr.Sint32Kind, pr.Sfixed32Kind:
		d.ValueTranscoder = &pgtype.Int4{Status: status}
		d.valueFunc = assertValueFunc(pr.ValueOfInt32)

	case pr.Int64Kind, pr.Sint64Kind, pr.Sfixed64Kind:
		d.ValueTranscoder = &pgtype.Int8{Status: status}
		d.valueFunc = assertValueFunc(pr.ValueOfInt64)

	case pr.FloatKind:
		d.ValueTranscoder = &pgtype.Float4{Status: status}
		d.valueFunc = assertValueFunc(pr.ValueOfFloat32)

	case pr.DoubleKind:
		d.ValueTranscoder = &pgtype.Float8{Status: status}
		d.valueFunc = assertValueFunc(pr.ValueOfFloat64)

	case pr.StringKind:
		d.ValueTranscoder = &pgtype.Text{Status: status}
		d.valueFunc = assertValueFunc(pr.ValueOfString)

	case pr.BytesKind:
		d.ValueTranscoder = &pgtype.Bytea{Status: status}
		d.valueFunc = assertValueFunc(pr.ValueOfBytes)

	case pr.Uint32Kind, pr.Fixed32Kind:
		d.ValueTranscoder = &pgtype.Int4{Status: status}
		d.valueFunc = convertIntValueFunc[int32](pr.ValueOfUint32)

	case pr.Uint64Kind, pr.Fixed64Kind:
		d.ValueTranscoder = &pgtype.Int8{Status: status}
		d.valueFunc = convertIntValueFunc[int64](pr.ValueOfUint64)

	case pr.MessageKind:
		name := fd.Message().FullName()

		switch name {
		case SupportedTimestamp:
			d.ValueTranscoder = &pgtype.Timestamptz{Status: status}
			d.valueFunc = timeStampValue

		default:
			return nil, fmt.Errorf("unsupported message type %q", name)
		}

	default:
		return nil, fmt.Errorf("unsupported type %q", fd.Kind())
	}

	return d, nil
}

const (
	SupportedTimestamp = "google.protobuf.Timestamp"
)

func timeStampValue(x interface{}) pr.Value {
	if x == nil {
		return pr.ValueOf(x)
	}

	return pr.ValueOfMessage(
		timestamppb.New(x.(time.Time)).ProtoReflect(),
	)
}
