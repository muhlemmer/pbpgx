package crud_test

import (
	"github.com/muhlemmer/pbpgx/crud"
	gen "github.com/muhlemmer/pbpgx/example_gen"
)

func ExampleNewTable() {
	columns := crud.Columns{"title": crud.Zero}
	crud.NewTable[gen.ProductColumns_Names, *gen.Product, int32]("public", "example", columns)
}
