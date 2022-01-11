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

package value

import (
	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

type listValue[T scalar] struct {
	pgtype.ValueTranscoder
	fd        pr.FieldDescriptor
	valueFunc func(T) pr.Value
}

func (v *listValue[T]) PGValue() pgtype.Value { return v.ValueTranscoder }

func (v *listValue[T]) SetTo(msg pr.Message) {
	var list []T
	v.AssignTo(&list)

	pl := msg.NewField(v.fd).List()
	for _, elem := range list {
		pl.Append(v.valueFunc(elem))
	}

	msg.Set(v.fd, pr.ValueOfList(pl))
}
