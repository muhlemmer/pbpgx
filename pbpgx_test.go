// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "embed"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	testCtx  context.Context
	errCtx   context.Context // allways errors
	connPool *pgxpool.Pool
	connDsn  = "user=pbpgx_tester host=db port=5432 dbname=pbpgx_tester sslmode=disable pool_max_conns=10"
)

func init() {
	if dsn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		connDsn = dsn
	}
}

func closeConnPool() {
	ctx, cancel := context.WithTimeout(testCtx, 5*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		connPool.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		panic(fmt.Errorf("closeConnPool: %w", ctx.Err()))
	}
}

func execQuerySlice(sqls string) {
	for _, sql := range strings.Split(sqls, ";") {
		func() {
			ctx, cancel := context.WithTimeout(testCtx, time.Second)
			defer cancel()

			if _, err := connPool.Exec(ctx, sql); err != nil {
				panic(err)
			}
		}()
	}
}

//go:embed testdata/create.sql
var createTablesSQL string

//go:embed testdata/drop.sql
var dropTablesSQL string

// testMain is a wrapper which allows defered statements to run before os.Exit is called.
func testMain(m *testing.M) int {
	var cancel context.CancelFunc
	testCtx, cancel = context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	errCtx, cancel = context.WithCancel(testCtx)
	cancel()

	pcfg, err := pgxpool.ParseConfig(connDsn)
	if err != nil {
		panic(err)
	}

	connPool, err = pgxpool.ConnectConfig(testCtx, pcfg)
	if err != nil {
		panic(err)
	}
	defer closeConnPool()

	execQuerySlice(dropTablesSQL)
	execQuerySlice(createTablesSQL)

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}
