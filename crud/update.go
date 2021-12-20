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
	"google.golang.org/protobuf/proto"
)

// UpdateOne record in a table, identified by id, with the contents of the data Message.
// Each field value will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
// Empty fields are omitted from the query and will not be modified.
func UpdateOne(ctx context.Context, x pbpgx.Executor, tab *Table, id interface{}, data proto.Message) (pgconn.CommandTag, error) {
	q := updateQuery(tab, data, tab.bufLens.updateOne.get(), ColumnWildcard, true)
	defer q.release()
	go tab.bufLens.updateOne.setHigher(q.Len())

	args, err := q.parseArgs(data, 1)
	if err != nil {
		return nil, fmt.Errorf("crud.UpdateOne: %w", err)
	}
	args = append(args, id)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return x.Exec(ctx, q.String(), args...)
}

// UpdateReturnOne updates one record in a Table, identified by id, with the contents of the data Message
// and returns the result in a message of type M.
// Each field value in req will be set to a corresponding column,
// matching on the procobuf fieldname, case sensitive.
// Empty fields are omitted from the query and will not be modified.
//
// The returned message will have the fields set as named by columns.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
// All returned columns must be present as field names in M.
// See scan for more details.
func UpdateReturnOne[M proto.Message, Col ColName](ctx context.Context, x pbpgx.Executor, tab *Table, id interface{}, data proto.Message, columns []Col) (msg M, err error) {
	q := updateQuery(tab, data, tab.bufLens.updateOne.get(), columns, true)
	defer q.release()
	go tab.bufLens.updateOne.setHigher(q.Len())

	args, err := q.parseArgs(data, 1)
	if err != nil {
		return msg, fmt.Errorf("crud.UpdateOne: %w", err)
	}
	args = append(args, id)

	return pbpgx.QueryRow[M](ctx, x, q.String(), args...)
}
