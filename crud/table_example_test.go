package crud_test

import (
	"github.com/muhlemmer/pbpgx/crud"
	gen "github.com/muhlemmer/pbpgx/example_gen"
	"github.com/muhlemmer/pbpgx/query"
)

func ExampleNewTable() {
	cd := query.ColumnDefault{"id": query.Zero}
	crud.NewTable[gen.ProductColumns_Names, *gen.Product, int32]("public", "example", cd)
}
