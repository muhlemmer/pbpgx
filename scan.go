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
