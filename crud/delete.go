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

	"github.com/jackc/pgconn"
	"github.com/muhlemmer/pbpgx"
	"google.golang.org/protobuf/proto"
)

// DeleteOne record from a Table, identified by id.
func DeleteOne(ctx context.Context, x pbpgx.Executor, tab *Table, id interface{}) (pgconn.CommandTag, error) {
	q := deleteQuery(tab, tab.bufLens.deleteOne.get(), noColumns)
	defer q.release()
	go tab.bufLens.deleteOne.setHigher(q.Len())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return x.Exec(ctx, q.String(), id)
}

// DeleteOne record from a Table, identified by id.
// Returns the deleted record in a message of type M.
//
// The returned message will have the fields set named by columns.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
// All columns must be present as field names in M.
// See scan for more details.
func DeleteReturnOne[M proto.Message, Col ColName](ctx context.Context, x pbpgx.Executor, tab *Table, id interface{}, columns []Col) (M, error) {
	q := deleteQuery(tab, tab.bufLens.deleteOne.get(), columns)
	defer q.release()
	go tab.bufLens.deleteOne.setHigher(q.Len())

	return pbpgx.QueryRow[M](ctx, x, q.String(), id)
}
