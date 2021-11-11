// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"constraints"
	"fmt"

	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
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

type destination struct {
	fieldDesc pr.FieldDescriptor
	valueDecoder
	valueFunc func(interface{}) pr.Value
}

func (d *destination) value() pr.Value {
	return d.valueFunc(d.Get())
}

func newDestination(fd pr.FieldDescriptor) *destination {
	d := &destination{
		fieldDesc: fd,
	}

	if fd.Cardinality() == pr.Repeated {
		panic(fmt.Errorf("unsupported type \"%s %s\" for scanning", pr.Repeated, fd.Kind()))
	}

	switch fd.Kind() {
	case pr.BoolKind:
		d.valueDecoder = new(pgtype.Bool)
		d.valueFunc = assertValueFunc(pr.ValueOfBool)

	case pr.Int32Kind, pr.Sint32Kind, pr.Sfixed32Kind:
		d.valueDecoder = new(pgtype.Int4)
		d.valueFunc = assertValueFunc(pr.ValueOfInt32)

	case pr.Int64Kind, pr.Sint64Kind, pr.Sfixed64Kind:
		d.valueDecoder = new(pgtype.Int8)
		d.valueFunc = assertValueFunc(pr.ValueOfInt64)

	case pr.FloatKind:
		d.valueDecoder = new(pgtype.Float4)
		d.valueFunc = assertValueFunc(pr.ValueOfFloat32)

	case pr.DoubleKind:
		d.valueDecoder = new(pgtype.Float8)
		d.valueFunc = assertValueFunc(pr.ValueOfFloat64)

	case pr.StringKind:
		d.valueDecoder = new(pgtype.Text)
		d.valueFunc = assertValueFunc(pr.ValueOfString)

	case pr.BytesKind:
		d.valueDecoder = new(pgtype.Bytea)
		d.valueFunc = assertValueFunc(pr.ValueOfBytes)

	case pr.Uint32Kind, pr.Fixed32Kind:
		d.valueDecoder = new(pgtype.Int4)
		d.valueFunc = convertIntValueFunc[int32](pr.ValueOfUint32)

	case pr.Uint64Kind, pr.Fixed64Kind:
		d.valueDecoder = new(pgtype.Int8)
		d.valueFunc = convertIntValueFunc[int64](pr.ValueOfUint64)

	default:
		panic(fmt.Errorf("unsupported type %q for scanning", fd.Kind()))
	}

	return d
}

func destinations(pfds pr.FieldDescriptors, pgfs []pgproto3.FieldDescription) []interface{} {
	fields := make([]interface{}, len(pgfs))

	for i, f := range pgfs {
		pfd := pfds.ByName(pr.Name(string(f.Name)))
		if pfd == nil {
			panic(fmt.Errorf("unknown field %s", f.Name))
		}

		fields[i] = newDestination(pfd)
	}

	return fields
}

func scanLimit[M proto.Message](limit int, rows pgx.Rows) (results []M, err error) {
	var m M

	defer func() {
		if err = checkNotRuntimeErr(recover()); err != nil {
			err = fmt.Errorf("pbpgx.Scan into proto.Message %T: %w", m, err)
		}
	}()

	msg := m.ProtoReflect()
	dest := destinations(msg.Descriptor().Fields(), rows.FieldDescriptions())

	for i := 0; rows.Next() && (limit == 0 || i < limit); i++ {
		msg = msg.New()

		if err := rows.Scan(dest...); err != nil {
			panic(err)
		}

		for _, v := range dest {
			d := v.(*destination)
			msg.Set(d.fieldDesc, d.value())
		}

		results = append(results, msg.Interface().(M))
	}

	return results, nil
}

// Scan returns a slice of proto messages of type M, filled with data from rows.
// It matches field names from rows to field names of the proto message type M.
// An error is returned if a column name in rows is not found in te message type's field names,
// if a matched message field is of an unsupported type or any scan error reported by the pgx driver.
func Scan[M proto.Message](rows pgx.Rows) ([]M, error) {
	return scanLimit[M](0, rows)
}
