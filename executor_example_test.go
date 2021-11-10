// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx_test

import (
	"context"

	"github.com/muhlemmer/pbpgx"
	gen "github.com/muhlemmer/pbpgx/example_gen"
)

func Query_Example() {
	var (
		result = new(gen.ProductList)
		err    error
	)

	result.Products, err = pbpgx.Query[*gen.Product](context.TODO(), conn, "select id, tile, price from products where id = $1;", 2)
	if err != nil {
		panic(err)
	}
}
