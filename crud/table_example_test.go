package crud_test

import (
	"github.com/muhlemmer/pbpgx/crud"
	gen "github.com/muhlemmer/pbpgx/example_gen"
)

func ExampleNewTable() {
	cd := crud.ColumnDefault{"id": crud.Zero}
	crud.NewTable[gen.ProductColumns_Names, *gen.Product, int32]("public", "example", cd)
}
