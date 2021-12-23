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
	"strings"

	"github.com/jackc/pgconn"
	"github.com/muhlemmer/pbpgx"
	"github.com/muhlemmer/pbpgx/query"
	"google.golang.org/protobuf/proto"
)

// CreateOne record in a Table, with the contents of the req Message.
// Each field value in data will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
// Empty fields are omitted from the query, the resulting values for the corresponding columns will depend on database defaults.
func (tab *Table[Col, Record, ID]) CreateOne(ctx context.Context, x pbpgx.Executor, data proto.Message) (pgconn.CommandTag, error) {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	fieldNames := b.Insert(tab.schema, tab.table, data, nil, true)

	args, err := query.ParseArgs(data, fieldNames, tab.cd)
	if err != nil {
		return nil, fmt.Errorf("crud.CreateOne: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return x.Exec(ctx, b.String(), args...)
}

// CreateReturnOne creates one record in a Table, with the contents of data
// and returns the result in a message of type Record.
// Each field value in data will be set to a corresponding column,
// matching on the procobuf fieldname, case sensitive.
// Empty fields are omitted from the query, the resulting values for the corresponding columns will depend on database defaults.
//
// The returned message will have the fields set as identified by returnColumns.
// If the length of columns is 0, the wildcard operator '*' is passed as columns spec to the database,
// returning all available columns in the table.
func (tab *Table[Col, Record, ID]) CreateReturnOne(ctx context.Context, x pbpgx.Executor, data proto.Message, returnColumns []Col) (m Record, err error) {
	b := tab.pool.Get()
	defer tab.pool.Put(b)

	fieldNames := b.Insert(tab.schema, tab.table, data, returnColumns, true)

	args, err := query.ParseArgs(data, fieldNames, tab.cd)
	if err != nil {
		var m Record
		return m, fmt.Errorf("crud.CreateOne: %w", err)
	}

	return pbpgx.QueryRow[Record](ctx, x, b.String(), args...)
}

// Create one or more records in a Table, with the contents of the req Message.
// Each field value in data will be set to a corresponding column,
// matching on the protobuf fieldname, case sensitive.
// Empty fields are omitted from the query, the resulting values for the corresponding columns will depend on database defaults.
func (tab *Table[Col, Record, ID]) Create(ctx context.Context, x Preparator, data []proto.Message) ([]pgconn.CommandTag, error) {
	if len(data) == 0 {
		return nil, nil
	}

	b := tab.pool.Get()
	defer tab.pool.Put(b)

	fieldNames := b.Insert(tab.schema, tab.table, data[0], nil, false)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sd, err := x.Prepare(ctx, strings.Join([]string{"CreateMany", tab.schema, tab.table}, "-"), b.String())
	if err != nil {
		return nil, fmt.Errorf("crud.CreateMany: %w", err)
	}

	tags := make([]pgconn.CommandTag, 0, len(data))

	for i, m := range data {
		args, err := query.ParseArgs(m, fieldNames, tab.cd)
		if err != nil {
			return tags, fmt.Errorf("crud.CreateMany[%d]: %w", i, err)
		}

		tag, err := x.Exec(ctx, sd.Name, args...)
		if err != nil {
			return tags, fmt.Errorf("crud.CreateMany[%d]: %w", i, err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
