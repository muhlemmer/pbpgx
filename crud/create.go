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
	"context"
	"fmt"

	"github.com/muhlemmer/pbpgx"
	"google.golang.org/protobuf/proto"
)

func (tab *Table[Col, Record, ID]) insertQuery(cols ColNames, returnColumns ...Col) (query string) {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	b.Insert(tab.schema, tab.table, cols, returnColumns...)

	return b.String()
}

// CreateOne creates one record in a Table, with the contents of data
// and returns the result in a message of type Record.
// Each field value in data will be set to a corresponding column from cols,
// matching on the procobuf fieldname, case sensitive.
//
// If any returnColumns are specified, the returned record will have the fields set as named by returnColumns.
// If no returnColumns, the returned record will always be nil.
func (tab *Table[Col, Record, ID]) CreateOne(ctx context.Context, x pbpgx.Executor, cols ColNames, data proto.Message, returnColumns []Col) (record Record, err error) {
	qs := tab.insertQuery(cols, returnColumns...)

	args, err := tab.columns.ParseArgs(data, cols)
	if err != nil {
		return record, fmt.Errorf("Table %s CreateOne: %w", tab.name(), err)
	}

	if len(returnColumns) > 0 {
		record, err = pbpgx.QueryRow[Record](ctx, x, qs, args...)
	} else {
		_, err = x.Exec(ctx, qs, args...)
	}

	if err != nil {
		return record, fmt.Errorf("Table %s CreateOne: %w", tab.name(), err)
	}

	return
}

// Create one or more records in a Table, with the contents of the req Message.
// Each field value in data will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
// Empty fields are omitted from the query, the resulting values for the corresponding columns will depend on database defaults.
//
// If any returnColumns are specified, the returned records will have the fields set as named by returnColumns.
// If no returnColumns, the returned slice will always be nil.
//
// Note: Although it is allowed to pass a connection pool as Executor,
// it is recommended to pass a Conn or Tx type.
// Create executes the same INSERT query for every entry in Data,
// And a connection pool does not reuse statements.
// If the connection is set up in "simple" mode, this function will likely have bad performance.
// Furthermore, this function makes no assumptions on transactional requirements of the call.
// Meaning that on error a part of the data may be inserted and will not be rolled back.
// It is the responsibilty of the caller to Begin and Rollback in case this is required.
func (tab *Table[Col, Record, ID]) Create(ctx context.Context, x pbpgx.Executor, cols ColNames, data []proto.Message, returnColumns ...Col) (records []Record, err error) {
	if len(cols) == 0 || len(data) == 0 {
		return nil, nil
	}

	qs := tab.insertQuery(cols, returnColumns...)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(returnColumns) > 0 {
		records = make([]Record, 0, len(data))
	}

	for i, m := range data {
		args, err := tab.columns.ParseArgs(m, cols)
		if err != nil {
			return records, fmt.Errorf("Table %s Create[%d]: %w", tab.name(), i, err)
		}

		var record Record

		if len(returnColumns) > 0 {
			record, err = pbpgx.QueryRow[Record](ctx, x, qs, args...)
		} else {
			_, err = x.Exec(ctx, qs, args...)
		}

		if err != nil {
			return records, fmt.Errorf("Table %s Create[%d]: %w", tab.name(), i, err)
		}

		if len(returnColumns) > 0 {
			records = append(records, record)
		}

	}

	return records, nil
}
