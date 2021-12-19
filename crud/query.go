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

package crud

import (
	"constraints"
	"fmt"
	"strconv"

	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/stringx"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

type query struct {
	*stringx.Builder
	tab        *Table
	fieldNames []string
	argPos     int
}

func (q *query) parseFields(req proto.Message, skipEmpty bool, exclude ...string) {
	exMap := make(map[string]struct{}, len(exclude))

	for _, s := range exclude {
		exMap[s] = struct{}{}
	}

	msg := req.ProtoReflect()
	fields := msg.Descriptor().Fields()
	l := fields.Len()

	q.fieldNames = make([]string, 0, l-len(exclude))

	for i := 0; i < l; i++ {
		fd := fields.Get(i)
		if skipEmpty && !msg.Has(fd) {
			continue
		}

		name := string(fd.Name())

		if _, ok := exMap[name]; !ok {
			q.fieldNames = append(q.fieldNames, name)
		}
	}
}

func (q *query) parseArgs(req proto.Message, extraCap int) (args []interface{}) {
	msg := req.ProtoReflect()
	fields := msg.Descriptor().Fields()

	args = make([]interface{}, 0, len(q.fieldNames)+extraCap)

	for _, name := range q.fieldNames {
		fd := fields.ByName(pr.Name(name))

		arg := pbpgx.NewValue(fd, q.tab.onEmpty[name].pgStatus())

		if msg.Has(fd) {
			arg.Set(msg.Get(fd).Interface())
		}

		args = append(args, arg)
	}

	return args
}

// writeIdentifier writes the full identifier of a table.
// This includes the schema, if it is set to Table.
// Written names will be quoted.
func (q *query) writeIdentifier() {
	if q.tab.schema != "" {
		q.WriteEnclosedString(q.tab.schema, stringx.DoubleQuotes)
		q.WriteByte(schemaSep)
	}

	q.WriteEnclosedString(q.tab.table, stringx.DoubleQuotes)
}

func (q *query) writePosArgs(n int) {
	for i := 0; i < n; i++ {
		if i != 0 {
			q.WriteString(columnSep)
		}

		q.argPos++

		q.WriteByte('$')
		q.WriteString(strconv.Itoa(q.argPos))
	}
}

func (q *query) release() {
	q.Reset()
	q.tab.pool.Put(q.Builder)
}

const (
	columnSep = ", "
	schemaSep = '.'
	idColumn  = "id"
	from      = " FROM "
	returning = " RETURNING "
	where     = " WHERE "
)

var NoRetCols = []fmt.Stringer{}

func writeColumns[Col fmt.Stringer](q *query, cols []Col) {
	for i, col := range cols {
		if i != 0 {
			q.WriteString(", ")
		}
		q.WriteEnclosedString(col.String(), stringx.DoubleQuotes)
	}
}

// INSERT INTO "public"."simple_rw" ("id", "title") VALUES ($1, $2) RETURNING "id";
func insertQuery[Col fmt.Stringer](tab *Table, req proto.Message, bufSize int, retCols []Col, skipEmpty bool, exclude ...string) *query {
	const (
		insertInto = "INSERT INTO "
		values     = " VALUES "
	)

	q := tab.newQuery(bufSize)
	q.parseFields(req, skipEmpty, exclude...)

	q.WriteString(insertInto)
	q.writeIdentifier()

	q.WriteString(" (")
	q.WriteEnclosedElements(q.fieldNames, columnSep, stringx.DoubleQuotes)
	q.WriteByte(')')

	q.WriteString(values)
	q.WriteByte('(')
	q.writePosArgs(len(q.fieldNames))
	q.WriteByte(')')

	if len(retCols) > 0 {
		q.WriteString(returning)
		writeColumns(q, retCols)
	}

	q.WriteByte(';')

	return q
}

// SELECT "id", "title", "price" FROM "public"."products" WHERE "id" = $1 LIMIT 1;
func selectQuery[ColDesc fmt.Stringer, I constraints.Integer](tab *Table, cols []ColDesc, bufSize int, wf whereFunc, limit I) *query {
	const (
		sselect = "SELECT "
		slimit  = " LIMIT "
	)

	q := tab.newQuery(bufSize)
	q.WriteString(sselect)

	if len(cols) > 0 {
		writeColumns(q, cols)
	} else {
		q.WriteByte('*')
	}

	q.WriteString(from)
	q.writeIdentifier()

	if wf != nil {
		wf(q)
	}

	if limit > 0 {
		q.WriteString(slimit)
		q.WriteString(strconv.FormatInt(int64(limit), 10))
	}

	q.WriteByte(';')

	return q
}

// UPDATE "public"."simple" SET "title" = $1, "data" = $2 WHERE "id" = $3 RETURNING ("title", "data");
func updateQuery[Col fmt.Stringer](tab *Table, data proto.Message, bufSize int, retCols []Col, skipEmpty bool, exclude ...string) *query {
	const (
		update = "UPDATE "
		set    = " SET "
		where  = " WHERE \"id\" = "
	)

	q := tab.newQuery(bufSize)
	q.parseFields(data, skipEmpty, exclude...)

	q.WriteString(update)
	q.writeIdentifier()

	q.WriteString(set)

	for i, name := range q.fieldNames {
		if i != 0 {
			q.WriteString(", ")
		}
		q.WriteEnclosedString(name, stringx.DoubleQuotes)
		q.WriteString(" = ")

		q.writePosArgs(1)
	}

	q.WriteString(where)

	q.writePosArgs(1)

	if len(retCols) > 0 {
		q.WriteString(returning)
		writeColumns(q, retCols)
	}

	q.WriteByte(';')
	return q
}

// DELETE FROM "public"."simple" WHERE "id" = $1 RETURNING ("id, "title", "data");
func deleteQuery[Col fmt.Stringer](tab *Table, bufSize int, retCols []Col) *query {
	const (
		delete = "DELETE"
		where  = " WHERE \"id\" = $1"
	)

	q := tab.newQuery(bufSize)

	q.WriteString(delete)
	q.WriteString(from)
	q.writeIdentifier()
	q.WriteString(where)

	if len(retCols) > 0 {
		q.WriteString(returning)
		writeColumns(q, retCols)
	}

	q.WriteByte(';')

	return q
}
