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

func (tab *Table[Col, Record, ID]) selectQuery(columns []Col, wf query.WhereFunc[Col], orderBy query.OrderWriter[Col], limit int64) string {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	b.Select(tab.schema, tab.table, columns, wf, orderBy, limit)
	return b.String()
}

// ReadOne record from a table, identified by id.
// The returned message will be of type Record,
// with the fields corresponding to columns populated.
func (tab *Table[Col, Record, ID]) ReadOne(ctx context.Context, x pbpgx.Executor, id ID, columns []Col) (record Record, err error) {
	record, err = pbpgx.QueryRow[Record](ctx, x, tab.selectQuery(columns, query.WhereID[Col], nil, 1), id)
	if err != nil {
		return record, fmt.Errorf("Table %s ReadOne: %w", tab.name(), err)
	}

	return record, nil
}

// ReadAll records up to limit from a table.
// The returned messages will be a slice of type Record,
// with the fields corresponding to columns populated.
func (tab *Table[Col, Record, ID]) ReadAll(ctx context.Context, x pbpgx.Executor, limit int64, columns []Col, orderBy query.OrderWriter[Col]) ([]Record, error) {
	records, err := pbpgx.Query[Record](ctx, x, tab.selectQuery(columns, nil, orderBy, limit))
	if err != nil {
		return nil, fmt.Errorf("Table %s ReadAll: %w", tab.name(), err)
	}

	return records, nil
}

// ReadList returns a list of records from a table, identified by ids.
// The returned messages will be a slice of type Record,
// with the fields corresponding to columns populated.
func (tab *Table[Col, Record, ID]) ReadList(ctx context.Context, x pbpgx.Executor, ids []ID, columns []Col, orderBy query.OrderWriter[Col]) ([]Record, error) {
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		args[i] = id
	}

	records, err := pbpgx.Query[Record](ctx, x, tab.selectQuery(columns, query.WhereIDInFunc[Col](len(ids)), orderBy, 0), args...)
	if err != nil {
		return nil, fmt.Errorf("Table %s ReadList: %w", tab.name(), err)
	}

	return records, nil
}
