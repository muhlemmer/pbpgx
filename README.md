[![Go](https://github.com/muhlemmer/pbpgx/actions/workflows/go.yml/badge.svg)](https://github.com/muhlemmer/pbpgx/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/muhlemmer/pbpgx/branch/main/graph/badge.svg?token=NS16O82CC1)](https://codecov.io/gh/muhlemmer/pbpgx)
[![Go Report Card](https://goreportcard.com/badge/github.com/muhlemmer/pbpgx)](https://goreportcard.com/report/github.com/muhlemmer/pbpgx)

# PBPGX

Package pbpgx provides a toolkit for easier Protocol Buffers interaction with PostgreSQL databases.

Pbpgx supports the Protocol Buffer types generated in your project, through protoreflect.
It is build on the pgx PostgreSQL driver toolkit for Go, for richer SQL type support.
Thefore pbpgx is targetting developers which are building protcol buffers based APIs against a PostgreSQL database.

This package is currently a WIP and is written against Go version 1.18-beta1 with type parameter support and
will support Go version 1.18 upward.

## Features

When not using an ORM a typical development workflow I often experience is:

1. Query preparation: Load definitions from files, execute templates or [build queries](#query-building) based on the request to the API.
2. [Query execution](#query-execution): Send the prepared query to the database, with metods such as `conn.Exec()` or `conn.Query()`.
3. [Row scanning](#row-scanning): In case of the read action `conn.Query()`, the returned data must be scanned with `rows.Next()` and `rows.Scan()`,
    checking for errors on each iteration.
    If the result columns are variable based on the request data, then the scan targets must also dynamically be prepared during query preparation.

Pbpgx can greatly reduce this effort for most boilerplate API tasks. This package can be used from the bottom-up.
For instance, if your project already has a way of preparing and executing queries, this package can still be used for [scanning](#row-scanning) Rows into protocol buffer messages.
However, for most [CRUD](https://en.wikipedia.org/wiki/CRUD) actions, this package provides a way to build and execute queries, returning the results in the required type in a single call.
It aims to minimize data transformation by the caller, by using Protocol Buffer messages as argument and return types.

### Row Scanning

A common design challange with Go is scanning database Rows results into structs. Ussualy, one would scan into the indivudual fields or local variables.
But what if the the selected columns are a variable in the application? In such cases developers have to resort to reflection, an ORM based on
reflection or a code generator like SQLBoiler with tons of adaptor code. 
Pbpgx uses [protoreflect](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect), which uses the embedded protobuf descriptors in the generated Go code.
This should give a potential performance edge over standard Go reflection. 

By using Go generics, the `Scan()` function is able to return slices of concrete message types, allowing you to consume the result set without further transformation. For example a product list in a [proto file](example_gen/example.proto):

```
message Product {
    int64 id = 1;
    string title = 2;
    double price = 3;
    google.protobuf.Timestamp created = 4;
}

message ProductList {
    repeated Product products = 1;
}

```

Can easily be filled with:

```
rows, err := conn.Query(ctx, "select id, title, price, created from products;")
if err != nil {
    panic(err)
}

result := new(gen.ProductList)

result.Products, err = pbpgx.Scan[*gen.Product](rows)
if err != nil {
    panic(err)
}

```

### Query Execution

The generic `Query()` function can be used to execute a query and return a slice of Protocol Buffer messages filled with the results, through [Row Scanning](#row-scanning):

```
result.Products, err = pbpgx.Query[*gen.Product](ctx, conn,
    "select title, price, created from products where id in ($1, $2, $3);", 2, 4, 5)
if err != nil {
    panic(err)
}
```

### CRUD operations

The `github.com/muhlemmer/pbpgx/crud` sub-package provides query building, execution and scanning for common CRUD operations. (Create, Read, Update and Delete).
This allows to easily integrate your database models with protocol buffers and derived RPC implementations, for all basic operations.
Argument and return types are typically of `proto.Message`, where applicable, again eliminating the need of adaptor code.
Also the use of generic type arguments and `protoreflect` ensures a maintainable, consice, code base.

#### Usage

Usage is simple. First we extend on our [protocal buffer definitions](example_gen/example.proto) from above:

```
message ProductColumns{
    enum Names {
        id = 0;
        title = 1;
        price = 2;
        created = 3;
    }
}

message ProductQuery {
    int64 id = 1;
    repeated ProductColumns.Names columns = 2;
}

message ProductCreateQuery {
    Product data = 1;
    repeated ProductColumns.Names columns = 2;
}
```

In our implementation code we reate a `Table` once, which takes care of query builder re-use:

```
columns := crud.Columns{"title": crud.Zero}
tab := crud.NewTable[gen.ProductColumns_Names, *gen.Product, int32]("public", "example", columns)
```

We can use the `id` and `columns` from an incomming query to read a single message,
returned as the requested type `Product`:

```
query := &gen.ProductQuery{
    Id: 2,
    Columns: []gen.ProductColumns_Names{
        gen.ProductColumns_title,
        gen.ProductColumns_price,
        gen.ProductColumns_created,
    },
}

record, err := tab.ReadOne(ctx, conn, query.Id, query.Columns)
if err != nil {
    panic(err)
}

```

Or, create a new entry (`INSERT`), returing specific fields after the database assigns the primary ID and created timestamp:

```
query := &gen.ProductCreateQuery{
    Data: &gen.Product{
        Title: "Great deal!",
        Price: 9.99,
    },
    Columns: []gen.ProductColumns_Names{
        gen.ProductColumns_id,
        gen.ProductColumns_title,
        gen.ProductColumns_price,
        gen.ProductColumns_created,
    },
}

record, err := tab.CreateOne(ctx, conn, crud.ParseFields(query.GetData(), true), query.GetData(), query.GetColumns()...)
if err != nil {
    panic(err)
}
```

Likewise, there are the `CreateOne`, `DeleteOne` and `UpdateOne` functions.

### Query building

The `crud` package is build on the `github.com/muhlemmer/pbpgx/query` sub-package. The building blocks for query building are exported, so that package consumers can use them to build custom queries.
The `query` package uses a `sync.Pool` for efficient memory re-use.
The builder's capacity is stored and on the next use, the same capacity is allocated in order to reduce the amount off allocations needed during building.
This results in a fairly efficient builder:

```
BenchmarkBuilder_Insert-16    	 5989922	       477.7 ns/op	     192 B/op	       1 allocs/op
BenchmarkBuilder_Select-16    	 2476291	       405.9 ns/op	      96 B/op	       1 allocs/op
BenchmarkBuilder_Update-16    	 3126738	       555.4 ns/op	     192 B/op	       1 allocs/op
BenchmarkBuilder_Delete-16    	 4884342	       303.8 ns/op	      80 B/op	       1 allocs/op
```

## License

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
