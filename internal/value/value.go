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
	"fmt"

	"github.com/jackc/pgtype"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

type Constructor func(pr.FieldDescriptor, pgtype.Status) (Value, error)

type register struct {
	messages map[pr.FullName]Constructor
}

func (r *register) addMessage(fullname pr.FullName, c Constructor) {
	if r.messages == nil {
		r.messages = map[pr.FullName]Constructor{}
	}

	r.messages[fullname] = c
}

func (r *register) newMessageValue(name pr.FullName, fd pr.FieldDescriptor, status pgtype.Status) (Value, error) {
	if c, ok := r.messages[name]; ok {
		return c(fd, status)
	}

	return nil, fmt.Errorf("value: message type %q not registered", name)
}

var registered register

func RegisterMessage(name pr.FullName, c Constructor) {
	registered.addMessage(name, c)
}

type Value interface {
	pgtype.ValueTranscoder
	PGValue() pgtype.Value
	SetTo(pr.Message)
}

func New(fd pr.FieldDescriptor, status pgtype.Status) (v Value, err error) {
	if msg := fd.Message(); msg != nil {
		return registered.newMessageValue(msg.FullName(), fd, status)
	}

	if fd.IsList() {
		return newlistValue(fd, status)
	}

	return newScalarValue(fd, status)
}
