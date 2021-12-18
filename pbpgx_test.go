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
	for i, sql := range strings.Split(sqls, ";") {
		func() {
			ctx, cancel := context.WithTimeout(testCtx, time.Second)
			defer cancel()

			if _, err := connPool.Exec(ctx, sql); err != nil {
				panic(fmt.Errorf("execQuerySlice[%d]: %w", i, err))
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
	testCtx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
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
