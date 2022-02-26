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

package crud

import (
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx/internal/value"
	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func contains[T constraints.Ordered](v T, list []T) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}

	return false
}

// ColNames expresses proto.Message field names for use as column specifiers.
type ColNames []string

// ParseFields returns a slice of field names from the passed proto message.
// Optionaly, empty (zero-value) fields can be skipped.
// Note that names are case sensitive, as defined in the proto file, not the Go struct field names.
func ParseFields(msg proto.Message, skipEmpty bool, ignore ...string) (cols ColNames) {
	rm := msg.ProtoReflect()
	fields := rm.Descriptor().Fields()

	cols = make([]string, 0, fields.Len()-len(ignore))

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		name := string(fd.Name())

		if skipEmpty && !rm.Has(fd) || contains(name, ignore) {
			continue
		}

		cols = append(cols, name)
	}

	return cols
}

type OnEmpty int

const (
	Null OnEmpty = iota // On empty field, set value to Null in the database.
	Zero                // On empty field, set value to its zero value.
)

func (o OnEmpty) pgStatus() pgtype.Status {
	if o == Zero {
		return pgtype.Present
	}

	return pgtype.Null
}

// Columns defines the default status of a column value.
// This affects the write behaviour of emtpy fields during Create (INSERT) and Update (UPDATE) calls.
type Columns map[string]OnEmpty

// ParseArgs parses the values from the fields named by fieldNames.
// The returned args contains pgtype values for efficient encoding.
// Empty fields will be set as `Null` by default, unless when set to `Zero`
// in Columns. Columns may be nil.
func (columns Columns) ParseArgs(msg proto.Message, colNames ColNames) (args []interface{}, err error) {
	rm := msg.ProtoReflect()
	fields := rm.Descriptor().Fields()

	args = make([]interface{}, 0, len(colNames)+5)

	for _, name := range colNames {
		fd := fields.ByName(pr.Name(name))
		if fd == nil {
			return nil, fmt.Errorf("ParseArgs: field %q not in msg %T", name, msg)
		}

		arg, err := value.New(fd, columns[name].pgStatus())
		if err != nil {
			return nil, fmt.Errorf("ParseArgs: %w", err)
		}

		if rm.Has(fd) {
			if fd.Kind() == pr.MessageKind {
				switch x := rm.Get(fd).Message().Interface().(type) {

				case *timestamppb.Timestamp:
					arg.Set(x.AsTime())

				default:
					// Ussualy a similar error is returned from NewValue.
					// Just checking here to be sure we didn't miss anything.
					return nil, fmt.Errorf("ParseArgs: unsupported message type %q", fd.Message().FullName())
				}

			} else {
				arg.Set(rm.Get(fd).Interface())
			}

		}

		args = append(args, arg)
	}

	return args, nil
}
