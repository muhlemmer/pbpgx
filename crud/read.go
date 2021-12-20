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
	"context"

	"github.com/muhlemmer/pbpgx"
	"google.golang.org/protobuf/proto"
)

// ReadOne record from a table, identified by id.
// The returned message will be of type M,
// with the fields corresponding to columns populated.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
// All columns must be present as field names in M.
// See scan for more details.
func ReadOne[M proto.Message, ID any, Col ColName](ctx context.Context, x pbpgx.Executor, tab *Table, id ID, columns []Col) (M, error) {
	q := selectQuery(tab, columns, tab.bufLens.read.n, whereID, 1)
	defer q.release()
	go tab.bufLens.createOne.setHigher(q.Len())

	return pbpgx.QueryRow[M](ctx, x, q.String(), id)
}

// ReadAll records up to limit from a table.
// The returned messages will be a slice of type M,
// with the fields corresponding to columns populated.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
// All columns must be present as field names in M.
// See scan for more details.
func ReadAll[M proto.Message, I constraints.Integer, Col ColName](ctx context.Context, x pbpgx.Executor, tab *Table, limit I, columns []Col) ([]M, error) {
	q := selectQuery(tab, columns, tab.bufLens.read.n, nil, limit)
	defer q.release()
	go tab.bufLens.createOne.setHigher(q.Len())

	return pbpgx.Query[M](ctx, x, q.String())
}

// ReadList returns a list of records from a table, identified by ids.
// The returned messages will be a slice of type M,
// with the fields corresponding to columns populated.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
// All columns must be present as field names in M.
// See scan for more details.
func ReadList[M proto.Message, ID any, Col ColName](ctx context.Context, x pbpgx.Executor, tab *Table, ids []ID, columns ...Col) ([]M, error) {
	q := selectQuery(tab, columns, tab.bufLens.read.n, whereIDInFunc(len(ids)), 0)
	defer q.release()
	go tab.bufLens.createOne.setHigher(q.Len())

	args := make([]interface{}, len(ids))

	for i, id := range ids {
		args[i] = id
	}

	return pbpgx.Query[M](ctx, x, q.String(), args...)
}
