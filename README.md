# PBPGX

Package pbpgx provides a toolkit for easier Protocol Buffers interaction with PostgreSQL databases.

Pbpgx supports the Protocol Buffer types generated in your project, through protoreflect.
It is build on the pgx PostgreSQL driver toolkit for Go, for richer SQL type support.
Thefore pbpgx is targetting developers which are building protcol buffers based APIs against a PostgreSQL database.

This package is currently a WIP and is written against the current `golang/go:master` branch with type parameter support and
will support Go version 1.18 upward.

## Features

### Scanning

A common design challange with Go is scanning database Rows results into structs. Ussualy, one would scan into the indivudual fields or local variables.
But what if the the selected columns are a variable in the application? In such cases developers have to resort to reflection, an ORM based on
reflection or a code generator like SQLBoiler with tons of adaptor code. 
Pbpgx uses protoreflect, which uses the embedded protobuf descriptors in the generated Go code.
This should give a potential performance edge over standard Go reflection. 

By using Go generics, the `Scan()` function is able to return slices of concrete message types, allowing you to consume the result set without further transformation. For example a product list in a [proto file](example_gen/example.proto):

```
message Product {
    int64 id = 1;
    string tile = 2;
    double price = 3;
}

message ProductList {
    repeated Product products = 1;
}

```

Can easily be filled with:

```
rows, err := conn.Query(context.Background(), "select id, tile, price from products;")
if err != nil {
    panic(err)
}

result := new(gen.ProductList)

result.Products, err = pbpgx.Scan[*gen.Product](rows)
if err != nil {
    panic(err)
}

```

## License

BSD 3-Clause "New" or "Revised" License
Copyright (c) 2021, Tim Möhlmann. All rights reserved.
Use of this source code is governed by a License that can be found in the [LICENSE](LICENSE) file.
