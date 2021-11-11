// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx_test

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/muhlemmer/pbpgx"
	gen "github.com/muhlemmer/pbpgx/example_gen"
	"google.golang.org/protobuf/encoding/protojson"
)

func ExampleQuery() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, "user=pbpgx_tester host=db port=5432 dbname=pbpgx_tester")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	result := new(gen.ProductList)

	result.Products, err = pbpgx.Query[*gen.Product](ctx, conn,
		"select title, price from products where id in ($1, $2, $3);", 2, 4, 5)
	if err != nil {
		panic(err)
	}

	out, _ := protojson.Marshal(result)
	fmt.Println(string(out))

	// Output: {"products":[{"title":"two","price":10.45},{"title":"four","price":100},{"title":"five","price":0.9}]}
}
