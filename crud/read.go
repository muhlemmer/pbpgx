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
	"google.golang.org/protobuf/proto"
)

// Selector is the proto.Message format accepted for Read query building
// and selecting single records after a write action.
type Selector[ID any, ColDesc fmt.Stringer] interface {
	proto.Message

	// GetId may return any type for ID.
	// The returned value will be used as-is as unique identifier in the "id"
	// column of the table being accessed.
	// The type must by compatible with the PostgreSQL type of the "id" column.
	// This method is ussualy generated for proto messages with an `id` field defined.
	GetId() ID

	// GetColumns returns a slice of fmt.Stringer, as column descriptors.
	// This method is ussualy generated for proto messages with a `columns` field defined.
	// This field should be a "repeated enum" type. Each enum value equals a column name.
	GetColumns() []ColDesc
}

// ReadOne record from a table, identified by query.Id.
// The returned message will be of type M,
// with the fields corresponding to query.Columns populated.
// All column names in query.Columns must be present as field names in M.
// See scan for more details.
func ReadOne[M proto.Message, ID any, ColDesc fmt.Stringer](ctx context.Context, x pbpgx.Executor, tab *Table, query Selector[ID, ColDesc]) (M, error) {
	q := selectQuery(tab, query.GetColumns(), tab.bufLens.read.n)
	defer q.release()
	go tab.bufLens.createOne.setHigher(q.Len())

	return pbpgx.QueryRow[M](ctx, x, q.String(), query.GetId())
}
