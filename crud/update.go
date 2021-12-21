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

	"github.com/jackc/pgconn"
	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/pbpgx/query"
	"google.golang.org/protobuf/proto"
)

// UpdateOne record in a table, identified by id, with the contents of the data Message.
// Each field value will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
// Empty fields are omitted from the query and will not be modified.
func (tab *Table[Col, Record, ID]) UpdateOne(ctx context.Context, x pbpgx.Executor, id ID, data proto.Message) (pgconn.CommandTag, error) {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	fieldNames := b.Update(tab.schema, tab.table, data, query.WhereID[Col], nil, true)

	args, err := query.ParseArgs(data, fieldNames, tab.cd)
	if err != nil {
		return nil, fmt.Errorf("crud.CreateOne: %w", err)
	}
	args = append(args, id)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return x.Exec(ctx, b.String(), args...)
}

// UpdateReturnOne updates one record in a Table, identified by id, with the contents of the data Message
// and returns the result in a message of type Record.
// Each field value in req will be set to a corresponding column,
// matching on the procobuf fieldname, case sensitive.
// Empty fields are omitted from the query and will not be modified.
//
// The returned message will have the fields set as named by returnColumns.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
func (tab *Table[Col, Record, ID]) UpdateReturnOne(ctx context.Context, x pbpgx.Executor, id ID, data proto.Message, returnColumns []Col) (msg Record, err error) {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	fieldNames := b.Update(tab.schema, tab.table, data, query.WhereID[Col], returnColumns, true)

	args, err := query.ParseArgs(data, fieldNames, tab.cd)
	if err != nil {
		var m Record
		return m, fmt.Errorf("crud.CreateOne: %w", err)
	}
	args = append(args, id)

	return pbpgx.QueryRow[Record](ctx, x, b.String(), args...)
}
