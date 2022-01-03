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
	"github.com/muhlemmer/pbpgx"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

// FieldNames expresses proto.Message field names for use as column specifiers.
type FieldNames []string

// ParseFields returns a slice of field names from the passed proto message.
// Optionaly, empty (zero-value) fields can be skipped.
// Note that names are case sensitive, as defined in the proto file, not the Go struct field names.
func ParseFields(msg proto.Message, skipEmpty bool) (fieldNames FieldNames) {
	rm := msg.ProtoReflect()
	fields := rm.Descriptor().Fields()

	fieldNames = make([]string, 0, fields.Len())

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if skipEmpty && !rm.Has(fd) {
			continue
		}

		fieldNames = append(fieldNames, string(fd.Name()))
	}

	return fieldNames
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

// ColumnDefault defines the default status of a column value.
// This affects the write behaviour of emtpy fields during Create (INSERT) and Update (UPDATE) calls.
type ColumnDefault map[string]OnEmpty

// ParseArgs parses the values from the fields named by fieldNames.
// The returned args contains pgtype values for efficient encoding.
// Empty fields will be set as `Null` by default, unless when set to `Zero`
// in the passed ColumnDefault. ColumnDefault may be nil.
func (fieldNames FieldNames) ParseArgs(msg proto.Message, cd ColumnDefault) (args []interface{}, err error) {
	rm := msg.ProtoReflect()
	fields := rm.Descriptor().Fields()

	args = make([]interface{}, 0, len(fieldNames)+5)

	for _, name := range fieldNames {
		fd := fields.ByName(pr.Name(name))

		arg, err := pbpgx.NewValue(fd, cd[name].pgStatus())
		if err != nil {
			return nil, fmt.Errorf("parseArgs: %w", err)
		}

		if rm.Has(fd) {
			arg.Set(rm.Get(fd).Interface())
		}

		args = append(args, arg)
	}

	return args, nil
}
