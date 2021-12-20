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

package pbpgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
)

// Executor is a common interface for database connection types.
// For example pgxpool.Pool, pgx[pool].Conn and pgx[pool].Tx all satisfy this interface.
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// Query runs the passed sql with args on the Executor x,
// and returns a slice of type M containing the results.
// See Scan for more details.
func Query[M proto.Message](ctx context.Context, x Executor, sql string, args ...interface{}) ([]M, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rows, err := x.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("pbpgx.Query: %w", err)
	}
	defer rows.Close()

	return Scan[M](rows)
}

// QueryRow runs the passed sql with args on the Executor x,
// and returns one row Scanned into a message of type M.
// See Scan for more details.
//
// In case of no rows, pgx.ErrNoRows is returned.
func QueryRow[M proto.Message](ctx context.Context, x Executor, sql string, args ...interface{}) (M, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rows, err := x.Query(ctx, sql, args...)
	if err != nil {
		var m M
		return m, fmt.Errorf("pbpgx.QueryRow: %w", err)
	}
	defer rows.Close()

	return ScanOne[M](rows)
}

// QueryStream runs the passed sql with args on the Executor x.
// Results are send to stream. See ScanStream for more details.
func QueryStream[M proto.Message](ctx context.Context, x Executor, stream ServerStream[M], sql string, args ...interface{}) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rows, err := x.Query(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("pbpgx.Query: %w", err)
	}
	defer rows.Close()

	return ScanStream(rows, stream)
}
