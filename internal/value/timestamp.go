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
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	SupportedTimestamp = "google.protobuf.Timestamp"
)

type timestampValue struct {
	pgtype.Timestamptz
	fd pr.FieldDescriptor
}

func (v *timestampValue) PGValue() pgtype.Value { return &v.Timestamptz }

func (v *timestampValue) SetTo(msg pr.Message) {
	msg.Set(v.fd, pr.ValueOfMessage(
		timestamppb.New(v.Time).ProtoReflect(),
	))
}

type timestampListValue struct {
	pgtype.TimestamptzArray
	fd pr.FieldDescriptor
}

func (v *timestampListValue) PGValue() pgtype.Value { return &v.TimestamptzArray }

func (v *timestampListValue) SetTo(msg pr.Message) {
	pl := msg.NewField(v.fd).List()

	for _, x := range v.Elements {
		pl.Append(pr.ValueOfMessage(
			timestamppb.New(x.Time).ProtoReflect(),
		))
	}

	msg.Set(v.fd, pr.ValueOfList(pl))
}
