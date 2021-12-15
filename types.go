// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"constraints"
	"fmt"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
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

type valueDecoder interface {
	pgtype.Value
	pgtype.BinaryDecoder
	pgtype.TextDecoder
}

type pgKind struct {
	pgtype.ValueTranscoder
	fieldDesc pr.FieldDescriptor
	valueFunc func(interface{}) pr.Value
}

func (d *pgKind) value() pr.Value {
	return d.valueFunc(d.Get())
}

func newPgKind(fd pr.FieldDescriptor) *pgKind {
	d := &pgKind{
		fieldDesc: fd,
	}

	if fd.Cardinality() == pr.Repeated {
		panic(fmt.Errorf("unsupported type \"%s %s\" for scanning", pr.Repeated, fd.Kind()))
	}

	switch fd.Kind() {
	case pr.BoolKind:
		d.ValueTranscoder = new(pgtype.Bool)
		d.valueFunc = assertValueFunc(pr.ValueOfBool)

	case pr.Int32Kind, pr.Sint32Kind, pr.Sfixed32Kind:
		d.ValueTranscoder = new(pgtype.Int4)
		d.valueFunc = assertValueFunc(pr.ValueOfInt32)

	case pr.Int64Kind, pr.Sint64Kind, pr.Sfixed64Kind:
		d.ValueTranscoder = new(pgtype.Int8)
		d.valueFunc = assertValueFunc(pr.ValueOfInt64)

	case pr.FloatKind:
		d.ValueTranscoder = new(pgtype.Float4)
		d.valueFunc = assertValueFunc(pr.ValueOfFloat32)

	case pr.DoubleKind:
		d.ValueTranscoder = new(pgtype.Float8)
		d.valueFunc = assertValueFunc(pr.ValueOfFloat64)

	case pr.StringKind:
		d.ValueTranscoder = new(pgtype.Text)
		d.valueFunc = assertValueFunc(pr.ValueOfString)

	case pr.BytesKind:
		d.ValueTranscoder = new(pgtype.Bytea)
		d.valueFunc = assertValueFunc(pr.ValueOfBytes)

	case pr.Uint32Kind, pr.Fixed32Kind:
		d.ValueTranscoder = new(pgtype.Int4)
		d.valueFunc = convertIntValueFunc[int32](pr.ValueOfUint32)

	case pr.Uint64Kind, pr.Fixed64Kind:
		d.ValueTranscoder = new(pgtype.Int8)
		d.valueFunc = convertIntValueFunc[int64](pr.ValueOfUint64)

	default:
		panic(fmt.Errorf("unsupported type %q for scanning", fd.Kind()))
	}

	return d
}
