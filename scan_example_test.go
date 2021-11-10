// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx_test

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/muhlemmer/pbpgx"
	gen "github.com/muhlemmer/pbpgx/example_gen"
)

var conn *pgx.Conn

func Scan_Example() {
	rows, err := conn.Query(context.TODO(), "select id, tile, price from products;")
	if err != nil {
		panic(err)
	}

	result := new(gen.ProductList)

	result.Products, err = pbpgx.Scan[*gen.Product](rows)
	if err != nil {
		panic(err)
	}
}
