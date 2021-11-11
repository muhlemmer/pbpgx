// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

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
}

func Exec(ctx context.Context, x Executor, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ct, err := x.Exec(ctx, sql, args...)
	if err != nil {
		return ct, fmt.Errorf("pbpgx.Exec: %w", err)
	}

	return ct, nil
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
