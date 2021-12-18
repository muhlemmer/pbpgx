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
	"context"
	"fmt"
	"strconv"

	"github.com/muhlemmer/stringx"
	"google.golang.org/protobuf/proto"
)

func readFields(req proto.Message) (n int, cols []string, posArgs []string, args []interface{}) {
	msg := req.ProtoReflect()
	fields := msg.Descriptor().Fields()
	l := fields.Len()

	cols = make([]string, 0, l)
	posArgs = make([]string, 0, l)
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

			arg := newPgKind(fd)
			arg.Set(msg.Get(fd).Interface())

			args = append(args, arg)
		}
	}

	return n, cols, posArgs, args
}

// INSERT INTO "public"."simple_rw" ("id", "title") VALUES ($1, $2) RETURNING "id";
func createOneQuery(schema, table string, req proto.Message) (n int, query string, args []interface{}) {
	const (
		insertInto = "INSERT INTO "
		values     = " VALUES "
		returning  = " RETURNING \"id\";"
	)

	n, cols, posArgs, args := readFields(req)

	if schema == "" {
		schema = DefaultSchema
	}

	n += len(insertInto) + len(schema) + len(table) + len(schemaSep) + 2*len(stringx.DoubleQuotes) + len(values) + len(returning) + len(columnSep)*(len(cols)-1)*2 + 5

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

	return n, b.String(), args
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

	_, query, args := createOneQuery(schema, table, req)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err = x.QueryRow(ctx, query, args...).Scan(&lastID); err != nil {
		panic(err)
	}

	return
}
