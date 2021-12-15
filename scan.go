// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

package pbpgx

import (
	"fmt"

	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

func destinations(pfds pr.FieldDescriptors, pgfs []pgproto3.FieldDescription) []interface{} {
	fields := make([]interface{}, len(pgfs))

	for i, f := range pgfs {
		pfd := pfds.ByName(pr.Name(string(f.Name)))
		if pfd == nil {
			panic(fmt.Errorf("unknown field %s", f.Name))
		}

		fields[i] = newPgKind(pfd)
	}

	return fields
}

func scanLimit[M proto.Message](limit int, rows pgx.Rows) (results []M, err error) {
	var m M

	defer func() {
		if err = checkNotRuntimeErr(recover()); err != nil {
			err = fmt.Errorf("pbpgx.Scan into proto.Message %T: %w", m, err)
		}
	}()

	msg := m.ProtoReflect()
	dest := destinations(msg.Descriptor().Fields(), rows.FieldDescriptions())

	for i := 0; rows.Next() && (limit == 0 || i < limit); i++ {
		msg = msg.New()

		if err := rows.Scan(dest...); err != nil {
			panic(err)
		}

		for _, v := range dest {
			d := v.(*pgKind)
			msg.Set(d.fieldDesc, d.value())
		}

		results = append(results, msg.Interface().(M))
	}

	return results, nil
}

// Scan returns a slice of proto messages of type M, filled with data from rows.
// It matches field names from rows to field names of the proto message type M.
// An error is returned if a column name in rows is not found in te message type's field names,
// if a matched message field is of an unsupported type or any scan error reported by the pgx driver.
func Scan[M proto.Message](rows pgx.Rows) ([]M, error) {
	return scanLimit[M](0, rows)
}
