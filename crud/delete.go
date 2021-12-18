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

// DeleteOne record from a Table, identified by id.
func DeleteOne(ctx context.Context, x pbpgx.Executor, tab *Table, id interface{}) (pgconn.CommandTag, error) {
	q := deleteQuery(tab, tab.bufLens.deleteOne.get(), NoRetCols)
	defer q.release()
	go tab.bufLens.deleteOne.setHigher(q.Len())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return x.Exec(ctx, q.String(), id)
}

// DeleteOne record from a Table, identified by selector.
// Returns the deleted record in a message of type M.
// The returned message will have the fields set as returned by selector.GetCols().
func DeleteReturnOne[M proto.Message, ID any, ColDesc fmt.Stringer](ctx context.Context, x pbpgx.Executor, tab *Table, selector Selector[ID, ColDesc]) (M, error) {
	q := deleteQuery(tab, tab.bufLens.deleteOne.get(), selector.GetColumns())
	defer q.release()
	go tab.bufLens.deleteOne.setHigher(q.Len())

	return pbpgx.QueryRow[M](ctx, x, q.String(), selector.GetId())
}
