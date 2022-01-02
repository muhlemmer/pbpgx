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

func (tab *Table[Col, Record, ID]) updateQuery(data proto.Message, wf query.WhereFunc[Col], skipEmpty bool, returnColumns ...Col) (qs string, fieldNames query.FieldNames) {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	fieldNames = query.ParseFields(data, skipEmpty)
	b.Update(tab.schema, tab.table, fieldNames, wf, returnColumns...)

	return b.String(), fieldNames
}

// UpdateOne updates one record in a Table, identified by id, with the contents of the data Message
// and returns the result in a message of type Record.
// Each field value in data will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
// Empty fields are omitted from the query and will not be modified.
//
// If any returnColumns are specified, the returned message will have the fields set as named by returnColumns.
// If no returnColumns, the returned message will always be nil.
func (tab *Table[Col, Record, ID]) UpdateOne(ctx context.Context, x pbpgx.Executor, id ID, data proto.Message, returnColumns ...Col) (record Record, err error) {
	qs, fieldNames := tab.updateQuery(data, query.WhereID[Col], true, returnColumns...)

	args, err := fieldNames.ParseArgs(data, tab.cd)
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
