# PBPGX

Package pbpgx provides a toolkit for easier Protocol Buffers interaction with PostgreSQL databases.

Pbpgx supports the Protocol Buffer types generated in your project, through protoreflect.
It is build on the pgx PostgreSQL driver toolkit for Go, for richer SQL type support.
Thefore pbpgx is targetting developers which are building protcol buffers based APIs against a PostgreSQL database.

This package is currently a WIP and is written against the current `golang/go:master` branch with type parameter support and
will support Go version 1.18 upward.

## Features

When not using an ORM a typical development workflow I often experience is:

1. Query definition: Either complete definitions, templates with partials or snipets in a custom query builder.
2. Query preparation: Load definitions from files, or execute templates or the query builder based on the request to the API.
3. [Query execution](#query-execution): Send the prepared query to the database, with metods such as `conn.Exec()` or `conn.Query()`.
4. [Row scanning](#row-scanning): In case of the read action `conn.Query()`, the returned data must be scanned with `rows.Next()` and `rows.Scan()`,
    checking for errors on each iteration.
    If the result columns are variable based on the request data, then the scan targets must also dynamically be prepared during query preparation.

Pbpgx can greatly reduce this effort for most boilerplate API tasks. This package can be used from the bottom-up.
For instance, if your project already has a way of preparing and executing queries, this package can still be used for [scanning](#row-scanning) Rows into protocol buffer messages.
However, for most [CRUD](https://en.wikipedia.org/wiki/CRUD) actions, this package provides a way to build and execute queries, returning the results in the required type in a single call.
It aims to minimize data transformation by the caller, by using Protocol Buffer messages as arguments and return types.

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
}

message ProductList {
    repeated Product products = 1;
}

```

Can easily be filled with:

```
rows, err := conn.Query(context.TODO(), "select id, tile, price from products;")
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
result.Products, err = pbpgx.Query[*gen.Product](context.TODO(), conn, "select id, tile, price from products where id = $1;", 2)
if err != nil {
    panic(err)
}
```

## License

BSD 3-Clause "New" or "Revised" License
Copyright (c) 2021, Tim Möhlmann. All rights reserved.
Use of this source code is governed by a License that can be found in the [LICENSE](LICENSE) file.
