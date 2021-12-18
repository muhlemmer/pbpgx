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
	"sync"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/stringx"
)

type bufLenCache struct {
	mu sync.RWMutex
	n  int
}

func (c *bufLenCache) setHigher(n int) {
	c.mu.Lock()
	if n > c.n {
		c.n = n
	}
	c.mu.Unlock()
}

func (c *bufLenCache) get() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.n
}

type OnEmpty int

const (
	Null OnEmpty = iota // On empty field, set value to Null in the database.
	Zero                // On empty field, set value to its zero value.
)

func (o OnEmpty) pgStatus() pgtype.Status {
	if o == Zero {
		return pgtype.Present
	}

	return pgtype.Null
}

// ColumnDefault defines the default status of a column value.
// This affects the write behaviour of emtpy fields during Create (INSERT) and Update (UPDATE) calls.
type ColumnDefault map[string]OnEmpty

// Table is a re-usable and concurrency safe object for query building.
// It holds reference to a table name, optional with schema
// and optimizes repeated query building for all supported CRUD functions in this package.
type Table struct {
	schema  string
	table   string
	onEmpty ColumnDefault

	pool sync.Pool

	bufLens struct {
		createOne bufLenCache
		read      bufLenCache
		updateOne bufLenCache
		deleteOne bufLenCache
	}
}

// NewTable returns a newly allocated table.
// Schema may be an empty string, in which case it will be ommitted from all queries built with this table.
// onEmpty specifies the behaviour when finding empty fields during data writes of multiple records.
// (protocol buffers does not have the notion of Null).
// Single record writes always ommit empty fields and ignore the onEmpty setting.
func NewTable(schema, table string, onEmpty ColumnDefault) *Table {
	return &Table{
		schema:  schema,
		table:   table,
		onEmpty: onEmpty,
	}
}

func (tab *Table) newQuery(size int) *query {
	b, ok := tab.pool.Get().(*stringx.Builder)
	if !ok {
		b = new(stringx.Builder)
	}

	b.Grow(size)

	return &query{
		Builder: b,
		tab:     tab,
	}
}
