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

package query

import (
	"fmt"
	"strconv"

	"github.com/muhlemmer/stringx"
)

const (
	columnSep = ", "
	schemaSep = '.'
	from      = " FROM "
)

// ColName is the generic interface for column name specification.
// Typically, one can use a protobuf enum for column names.
// The enum entries have a String method, implementing this interface.
type ColName interface {
	fmt.Stringer
}

// Builder for queries
type Builder[Col ColName] struct {
	stringx.Builder
	argPos  int
	lastCap int
}

func (b *Builder[Col]) Reset() {
	b.argPos = 0
	b.lastCap = b.Cap()
	b.Builder.Reset()
}

// WriteIdentifier writes the full identifier of a table.
// This includes the schema, if it is not empty.
// Written names will be quoted.
func (b *Builder[Col]) WriteIdentifier(schema, table string) {
	if schema != "" {
		b.WriteEnclosedString(schema, stringx.DoubleQuotes)
		b.WriteByte(schemaSep)
	}

	b.WriteEnclosedString(table, stringx.DoubleQuotes)
}

// WritePosArgs writes n amount of positional arguments to the query,
// in the form of "$1, $2".
// WritePosArgs can be called multiple times while building a query,
// and all position values in the resulting query will be uniquely incremented values.
func (b *Builder[Col]) WritePosArgs(n int) {
	for i := 0; i < n; i++ {
		if i != 0 {
			b.WriteString(columnSep)
		}

		b.argPos++

		b.WriteByte('$')
		b.WriteString(strconv.Itoa(b.argPos))
	}
}

// WriteColumnSpec writes the column specifier to the query.
// All column names are doube-quoted and comma seperated.
func (b *Builder[Col]) WriteColumnSpec(columns []Col) {
	for i, col := range columns {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteEnclosedString(col.String(), stringx.DoubleQuotes)
	}
}

// WriteReturnClause writes a RETURNING clause with a column spec.
// If columns has zero length, nothing will be written (no-op).
func (b *Builder[Col]) WriteReturnClause(columns []Col) {
	const returning = " RETURNING "

	if len(columns) > 0 {
		b.WriteString(returning)
		b.WriteColumnSpec(columns)
	}
}

// Insert builds an insert query.
// See WriteReturnClause on when and how the RETURNING clause is written.
//   INSERT INTO "public"."simple_rw" ("id", "title") VALUES ($1, $2) RETURNING "id";
func (b *Builder[Col]) Insert(schema, table string, insertColumns FieldNames, returnColumns ...Col) {
	const (
		insertInto = "INSERT INTO "
		values     = " VALUES "
	)

	b.WriteString(insertInto)
	b.WriteIdentifier(schema, table)

	b.WriteString(" (")
	b.WriteEnclosedElements(insertColumns, columnSep, stringx.DoubleQuotes)
	b.WriteByte(')')

	b.WriteString(values)
	b.WriteByte('(')
	b.WritePosArgs(len(insertColumns))
	b.WriteByte(')')

	b.WriteReturnClause(returnColumns)

	b.WriteByte(';')
}

// Select builds a select query.
// The WHERE clause must be written by the passed WhereFunc, which will not be called if nil.
// The LIMIT clause is only written when greater than 0.
//   SELECT "id", "title", "price" FROM "public"."products" WHERE "id" = $1 LIMIT 1;
func (b *Builder[Col]) Select(schema, table string, columns []Col, wf WhereFunc[Col], limit int64) {
	const (
		sselect = "SELECT "
		slimit  = " LIMIT "
	)

	b.WriteString(sselect)

	b.WriteColumnSpec(columns)

	b.WriteString(from)
	b.WriteIdentifier(schema, table)

	if wf != nil {
		wf(b)
	}

	if limit > 0 {
		b.WriteString(slimit)
		b.WriteString(strconv.FormatInt(limit, 10))
	}

	b.WriteByte(';')
}

// Update builds a update query.
// The WHERE clause must be written by the passed WhereFunc, which will not be called if nil.
// WARNING: a nil WhereFunc will result in updates of all records in a table!
// See WriteReturnClause on when and how the RETURNING clause is written.
//   UPDATE "public"."simple" SET "title" = $1, "data" = $2 WHERE "id" = $3 RETURNING ("title", "data");
func (b *Builder[Col]) Update(schema, table string, updateColumns FieldNames, wf WhereFunc[Col], returnColumns ...Col) {
	const (
		update = "UPDATE "
		set    = " SET "
	)

	b.WriteString(update)
	b.WriteIdentifier(schema, table)

	b.WriteString(set)

	for i, name := range updateColumns {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteEnclosedString(name, stringx.DoubleQuotes)
		b.WriteString(" = ")

		b.WritePosArgs(1)
	}

	if wf != nil {
		wf(b)
	}

	b.WriteReturnClause(returnColumns)
	b.WriteByte(';')

}

// Delete builds a delete query.
// The WHERE clause must be written by the passed WhereFunc, which will not be called if nil.
// WARNING: a nil WhereFunc will result in deletion of all records in a table!
// See WriteReturnClause on when and how the RETURNING clause is written.
//   DELETE FROM "public"."simple" WHERE "id" = $1 RETURNING ("id, "title", "data");
func (b *Builder[Col]) Delete(schema, table string, wf WhereFunc[Col], returnColumns ...Col) {
	const (
		delete = "DELETE"
	)

	b.WriteString(delete)
	b.WriteString(from)
	b.WriteIdentifier(schema, table)

	if wf != nil {
		wf(b)
	}

	b.WriteReturnClause(returnColumns)
	b.WriteByte(';')
}
