// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/stringx"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

func newArg(msg pr.Message, fd pr.FieldDescriptor) (v pgtype.Value) {
	if fd.Cardinality() == pr.Repeated {
		panic(fmt.Errorf("unsupported type \"%s %s\" for args", pr.Repeated, fd.Kind()))
	}

	switch fd.Kind() {
	case pr.BoolKind:
		v = new(pgtype.Bool)

	case pr.Int32Kind, pr.Sint32Kind, pr.Sfixed32Kind:
		v = new(pgtype.Int4)

	case pr.Int64Kind, pr.Sint64Kind, pr.Sfixed64Kind:
		v = new(pgtype.Int8)

	case pr.FloatKind:
		v = new(pgtype.Float4)

	case pr.DoubleKind:
		v = new(pgtype.Float8)

	case pr.StringKind:
		v = new(pgtype.Text)

	case pr.BytesKind:
		v = new(pgtype.Bytea)

	case pr.Uint32Kind, pr.Fixed32Kind:
		v = new(pgtype.Int4)

	case pr.Uint64Kind, pr.Fixed64Kind:
		v = new(pgtype.Int8)

	default:
		panic(fmt.Errorf("unsupported type %q for args", fd.Kind()))
	}

	v.Set(msg.Get(fd).Interface())

	return v
}

// INSERT INTO "public"."simple" ("id", "title") VALUES ($1, $2) RETURNING "id";
func createOneQuery(schema, table string, req proto.Message) (query string, n int, args []interface{}) {
	const (
		insertInto = "INSERT INTO "
		values     = " VALUES "
		returning  = " RETURNING \"id\";"
	)

	if schema == "" {
		schema = DefaultSchema
	}

	n += len(insertInto) + len(schema) + len(table) + len(schemaSep) + 2*len(stringx.DoubleQuotes) + len(values) + len(returning)

	msg := req.ProtoReflect()
	fields := msg.Descriptor().Fields()
	l := fields.Len()

	cols := make([]string, 0, l)
	posArgs := make([]string, 0, l)
	args = make([]interface{}, 0, l)

	var pos int

	for i := 0; i < l; i++ {
		fd := fields.Get(i)
		if msg.Has(fd) {
			name := fd.Name()
			cols = append(cols, string(name))
			n += len(name) + len(stringx.DoubleQuotes)

			pos++
			pa := "$" + strconv.Itoa(pos)
			posArgs = append(posArgs, pa)
			n += len(pa)

			args = append(args, newArg(msg, fd))
		}

	}

	n += len(columnSep)*(len(cols)-1)*2 + 5

	var b stringx.Builder
	b.Grow(n)

	b.WriteString(insertInto)
	b.WriteEnclosedElements([]string{schema, table}, schemaSep, stringx.DoubleQuotes)
	b.WriteString(" (")
	b.WriteEnclosedElements(cols, columnSep, stringx.DoubleQuotes)
	b.WriteByte(')')

	b.WriteString(values)
	b.WriteEnclosedJoin(posArgs, columnSep, stringx.RoundBrackets)
	b.WriteString(returning)

	return b.String(), n, args
}

// CreateOne record, with the contents of req.
// Empty fields will be ommited from the the INSERT query, and their actual value will depend
// on the column default, which might be NULL or other default value definition.
// Fields much match name and type of the table columns.
// lastID is returned with the value of the "id" column for the newly inserted record.
func CreateOne[M proto.Message](ctx context.Context, x Executor, schema, table string, req proto.Message) (lastID interface{}, err error) {
	defer func() {
		if err = checkNotRuntimeErr(recover()); err != nil {
			err = fmt.Errorf("pbpgx.CreateOne with proto.Message %T: %w", req, err)
		}
	}()

	query, _, args := createOneQuery(schema, table, req)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err = x.QueryRow(ctx, query, args...).Scan(&lastID); err != nil {
		panic(err)
	}

	return
}
