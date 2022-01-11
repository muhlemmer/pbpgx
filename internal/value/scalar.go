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
	"constraints"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

type signedScalar interface {
	int32 | int64
}

type unsignedScalar interface {
	uint32 | uint64
}

type scalar interface {
	signedScalar | unsignedScalar | bool | float32 | float64 | string | []byte
}

type scalarValue[T scalar] struct {
	pgtype.ValueTranscoder
	fd        pr.FieldDescriptor
	valueFunc func(T) pr.Value
}

func (v *scalarValue[T]) PGValue() pgtype.Value { return v.ValueTranscoder }

func (v *scalarValue[T]) SetTo(msg pr.Message) {
	if pgv, ok := v.Get().(T); ok {
		msg.Set(v.fd, v.valueFunc(pgv))
	}
}

// convertIntValueFunc returns a function which calls passed Value function
// with value of A converted the B.
func convertIntValueFunc[A, B constraints.Integer](f func(B) pr.Value) func(a A) pr.Value {
	return func(a A) pr.Value {
		return f(B(a))
	}
}
