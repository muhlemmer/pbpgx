// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"context"
	"fmt"

	"github.com/muhlemmer/stringx"
	"google.golang.org/protobuf/proto"
)

const (
	columnSep = ", "
	schemaSep = "."
)

func colSlice[ColType fmt.Stringer](columns []ColType) ([]string, int) {
	if len(columns) == 0 {
		return nil, 1
	}

	n := len(columnSep) * (len(columns) - 1)

	out := make([]string, len(columns))
	for i, s := range columns {
		out[i] = s.String()
		n += len(out[i]) + len(stringx.DoubleQuotes)
	}

	return out, n
}

var (
	DefaultSchema = "public"
)

// SELECT "id", "title", "price" FROM "public"."products" WHERE "id" = $1 LIMIT 1;
func readOneQuery[ColDesc fmt.Stringer](schema, table string, columns []ColDesc) (string, int) {
	const (
		sselect    = "SELECT "
		from       = " FROM "
		whereLimit = " WHERE \"id\" = $1 LIMIT 1;"
	)

	if schema == "" {
		schema = DefaultSchema
	}

	cols, n := colSlice(columns)
	n += len(sselect) + len(from) + len(schemaSep) + len(whereLimit) + len(schema) + len(table) + 2*len(stringx.DoubleQuotes)

	var b stringx.Builder
	b.Grow(n)

	b.WriteString(sselect)

	if len(cols) > 0 {
		b.WriteEnclosedElements(cols, columnSep, stringx.DoubleQuotes)
	} else {
		b.WriteByte('*')
	}

	b.WriteString(from)
	b.WriteEnclosedElements([]string{schema, table}, schemaSep, stringx.DoubleQuotes)
	b.WriteString(whereLimit)

	return b.String(), n
}

// ReadOneQuery is the proto.Message format accepted for ReadOne query building.
type ReadOneQuery[ID any, ColDesc fmt.Stringer] interface {
	proto.Message

	// GetId may return any type for ID.
	// The returned value will be used as-is as unique identifier in the "id"
	// column of the table being accessed.
	// The type must by compatible with the PostgreSQL type of the "id" column.
	// This method is ussualy generated for proto messages with an `id` field defined.
	GetId() ID

	// GetColumns returns a slice of fmt.Stringer, as column descriptors.
	// This method is ussualy generated for proto messages with a `columns` field defined.
	// This field should be a "repeated enum" type. Each enum value equals a column name.
	GetColumns() []ColDesc
}

// ReadOne record, identified by query.Id.
// The returned message will be of type M,
// with the fields corresponding to query.Columns populated.
// All column names in query.Columns must be present as field names in M.
// See scan for more details.
func ReadOne[ID any, ColDesc fmt.Stringer, M proto.Message](ctx context.Context, x Executor, schema, table string, query ReadOneQuery[ID, ColDesc]) (M, error) {
	qs, _ := readOneQuery(schema, table, query.GetColumns())
	return QueryRow[M](ctx, x, qs, query.GetId())
}
