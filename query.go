// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
)

type Executor interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type QueryArgs interface {
	Query() string
	Args() []interface{}
}

func Query[M proto.Message](ctx context.Context, x Executor, qa QueryArgs) ([]M, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rows, err := x.Query(ctx, qa.Query(), qa.Args()...)
	if err != nil {
		return nil, fmt.Errorf("pbpgx.Query: %w", err)
	}
	defer rows.Close()

	return Scan[M](rows)
}
