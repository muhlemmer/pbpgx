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
func UpdateOne[ID any, ColDesc fmt.Stringer](ctx context.Context, x pbpgx.Executor, tab *Table, id interface{}, data proto.Message) (pgconn.CommandTag, error) {
	q := updateQuery(tab, data, tab.bufLens.updateOne.get(), NoRetCols, true)
	defer q.release()
	go tab.bufLens.updateOne.setHigher(q.Len())

	args := append(q.parseArgs(data, 1), id)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return x.Exec(ctx, q.String(), args...)
}

// UpdateReturnOne updates one record in a Table, identified by selector, with the contents of the data Message
// and returns the result in a message of type M.
// Each field value in req will be set to a corresponding column,
// matching on the procobuf fieldname, case sensitive.
// Empty fields are omitted from the query and will not be modified.
// The returned message will have the fields set as returned by selector.GetColumns().
func UpdateReturnOne[M proto.Message, ID any, ColDesc fmt.Stringer](ctx context.Context, x pbpgx.Executor, tab *Table, selector Selector[ID, ColDesc], data proto.Message) (M, error) {
	q := updateQuery(tab, data, tab.bufLens.updateOne.get(), selector.GetColumns(), true)
	defer q.release()
	go tab.bufLens.updateOne.setHigher(q.Len())

	args := append(q.parseArgs(data, 1), selector.GetId())
	return pbpgx.QueryRow[M](ctx, x, q.String(), args...)
}
