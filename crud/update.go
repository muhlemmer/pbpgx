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
	"github.com/muhlemmer/pbpgx/query"
	"google.golang.org/protobuf/proto"
)

func (tab *Table[Col, Record, ID]) updateQuery(cols ColNames, wf query.WhereFunc[Col], returnColumns ...Col) (qs string) {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	b.Update(tab.schema, tab.table, cols, wf, returnColumns...)

	return b.String()
}

// UpdateOne updates one record in a Table, identified by id, with the contents of the data Message
// and returns the result in a message of type Record.
// Each field value in data will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
//
// If any returnColumns are specified, the returned message will have the fields set as named by returnColumns.
// If no returnColumns, the returned message will always be nil.
func (tab *Table[Col, Record, ID]) UpdateOne(ctx context.Context, x pbpgx.Executor, cols ColNames, id ID, data proto.Message, returnColumns ...Col) (record Record, err error) {
	qs := tab.updateQuery(cols, query.WhereID[Col], returnColumns...)

	args, err := tab.columns.ParseArgs(data, cols)
	if err != nil {
		return record, fmt.Errorf("Table %s UpdateOne: %w", tab.name(), err)
	}
	args = append(args, id)

	if len(returnColumns) > 0 {
		record, err = pbpgx.QueryRow[Record](ctx, x, qs, args...)
	} else {
		_, err = x.Exec(ctx, qs, args...)
	}

	if err != nil {
		return record, fmt.Errorf("Table %s UpdateOne: %w", tab.name(), err)
	}

	return record, err
}
