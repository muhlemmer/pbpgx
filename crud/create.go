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

// CreateOne record in a Table, with the contents of the req Message.
// Each field value will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
// Empty fields are omitted from the query, the resulting values for the corresponding columns will depend on database defaults.
func CreateOne(ctx context.Context, x pbpgx.Executor, tab *Table, req proto.Message) (pgconn.CommandTag, error) {
	q := insertQuery(tab, req, tab.bufLens.createOne.get(), noColumns, true)
	defer q.release()
	go tab.bufLens.createOne.setHigher(q.Len())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	args, err := q.parseArgs(req, 0)
	if err != nil {
		return nil, fmt.Errorf("crud.CreateOne: %w", err)
	}

	return x.Exec(ctx, q.String(), args...)
}

// CreateReturnOne creates one record in a Table, with the contents of the req Message
// and returns the result in a message of type M.
// Each field value in req will be set to a corresponding column,
// matching on the procobuf fieldname, case sensitive.
// Empty fields are omitted from the query, the resulting values for the corresponding columns will depend on database defaults.
//
// The returned message will have the fields set as identified by columns.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
// All returned columns must be present as field names in M.
// See scan for more details.
func CreateReturnOne[M proto.Message, ColDesc fmt.Stringer](ctx context.Context, x pbpgx.Executor, tab *Table, req proto.Message, columns []ColDesc) (m M, err error) {
	q := insertQuery(tab, req, tab.bufLens.createOne.get(), columns, true)
	defer q.release()
	go tab.bufLens.createOne.setHigher(q.Len())

	args, err := q.parseArgs(req, 0)
	if err != nil {
		return m, fmt.Errorf("crud.CreateOne: %w", err)
	}

	return pbpgx.QueryRow[M](ctx, x, q.String(), args...)
}
