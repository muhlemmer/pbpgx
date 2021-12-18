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

func ExampleScan() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, "user=pbpgx_tester host=db port=5432 dbname=pbpgx_tester")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "select id, title, price from products;")
	if err != nil {
		panic(err)
	}

	result := new(gen.ProductList)

	result.Products, err = pbpgx.Scan[*gen.Product](rows)
	if err != nil {
		panic(err)
	}

	out, _ := protojson.Marshal(result)
	fmt.Println(string(out))

	// {"products":[{"id":"1","title":"one","price":9.99},{"id":"2","title":"two","price":10.45},{"id":"3","title":"three"},{"id":"4","title":"four","price":100},{"id":"5","title":"five","price":0.9}]}
}
