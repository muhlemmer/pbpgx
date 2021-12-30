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

	"github.com/muhlemmer/pbpgx/query"
	"google.golang.org/protobuf/proto"
)

// Table is a re-usable and concurrency safe object for CRUD operations.
// It holds reference to a table name, optional with schema
// and optimizes repeated query building for all supported CRUD functions in this package.
type Table[Col query.ColName, Record proto.Message, ID constraints.Ordered] struct {
	schema string
	table  string
	cd     query.ColumnDefault
	pool   query.Pool[Col]
}

// NewTable returns a newly allocated table.
// Schema may be an empty string, in which case it will be ommitted from all queries built for this table.
// ColumnDefault specifies the behaviour when finding empty fields during data writes of multiple records.
// (protocol buffers does not have the notion of Null).
// Single record writes always ommit empty fields and ignore the onEmpty setting.
//
// The Col type parameter is typically an entry of a protocol buffers enum with columns names. (must implement the String() method).
// The Record type parameter should be a protocol buffer message representing the databae schema.
// Column names must match with field names, case sensitive.
// It is recommended to define a field for each column, for usage with the wildcard operator '*'.
// It is safe to have more fields than columns, the surplus will be ignored.
// See pbpgx.Scan for details.
// The ID type parameter should match the type used in "id" column of the table,
// used to match a single row in Read, Update and Delete.
func NewTable[Col query.ColName, Record proto.Message, ID constraints.Ordered](schema, table string, cd query.ColumnDefault) *Table[Col, Record, ID] {
	return &Table[Col, Record, ID]{
		schema: schema,
		table:  table,
		cd:     cd,
	}
}

func (tab *Table[Col, Record, ID]) name() string {
	var b query.Builder[Col]
	b.WriteIdentifier(tab.schema, tab.table)

	return b.String()
}
