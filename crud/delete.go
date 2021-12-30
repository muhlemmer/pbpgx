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
)

func (tab *Table[Col, Record, ID]) deleteQuery(wf query.WhereFunc[Col], returnColumns ...Col) string {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	b.Delete(tab.schema, tab.table, query.WhereID[Col], returnColumns)
	return b.String()
}

// DeleteOne record from a Table, identified by id.
// Returns the deleted record in a message of type Record.
//
// If any returnColumns are specified, the returned message will have the fields set as named by returnColumns.
// If no returnColumns, the returned message will always be nil.
func (tab *Table[Col, Record, ID]) DeleteOne(ctx context.Context, x pbpgx.Executor, id ID, returnColumns ...Col) (record Record, err error) {
	qs := tab.deleteQuery(query.WhereID[Col], returnColumns...)

	if len(returnColumns) > 0 {
		record, err = pbpgx.QueryRow[Record](ctx, x, qs, id)
	} else {
		_, err = x.Exec(ctx, qs, id)
	}

	if err != nil {
		return record, fmt.Errorf("Table %s DeleteOne: %w", tab.name(), err)
	}

	return record, nil
}
