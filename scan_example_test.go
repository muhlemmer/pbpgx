package pbpgx_test

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/muhlemmer/pbpgx"
	gen "github.com/muhlemmer/pbpgx/example_gen"
)

var conn *pgx.Conn

func Scan_Example() {
	rows, err := conn.Query(context.Background(), "select id, tile, price from products;")
	if err != nil {
		panic(err)
	}

	result := new(gen.ProductList)

	result.Products, err = pbpgx.Scan[*gen.Product](rows)
	if err != nil {
		panic(err)
	}
}
