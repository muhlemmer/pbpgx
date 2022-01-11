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
	"context"
	"fmt"

	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
	pr "google.golang.org/protobuf/reflect/protoreflect"
)

func destinations(pfds pr.FieldDescriptors, pgfs []pgproto3.FieldDescription) ([]interface{}, error) {
	fields := make([]interface{}, len(pgfs))

	for i, f := range pgfs {
		pfd := pfds.ByName(pr.Name(string(f.Name)))
		if pfd == nil {
			return nil, fmt.Errorf("unknown field %s", f.Name)
		}

		v, err := NewValue(pfd, pgtype.Undefined)
		if err != nil {
			return nil, err
		}

		fields[i] = v
	}

	return fields, nil
}

type scanner[M proto.Message] struct {
	rows pgx.Rows
	msg  pr.Message
	dest []interface{}
}

func newScanner[M proto.Message](rows pgx.Rows) (*scanner[M], error) {
	var m M

	msg := m.ProtoReflect()
	dest, err := destinations(msg.Descriptor().Fields(), rows.FieldDescriptions())
	if err != nil {
		return nil, err
	}

	return &scanner[M]{
		rows: rows,
		msg:  msg,
		dest: dest,
	}, nil
}

func (s *scanner[M]) scanRow() (M, error) {
	msg := s.msg.New()

	if err := s.rows.Scan(s.dest...); err != nil {
		var m M
		return m, fmt.Errorf("pbpgx.Scan into proto.Message %T: %w", m, err)
	}

	for _, d := range s.dest {
		d.(Value).setTo(msg)
	}

	return msg.Interface().(M), nil
}

// Scan returns a slice of proto messages of type M, filled with data from rows.
// It matches field names from rows to field names of the proto message type M.
// An error is returned if a column name in rows is not found in te message type's field names,
// if a matched message field is of an unsupported type or any scan error reported by the pgx driver.
func Scan[M proto.Message](rows pgx.Rows) (result []M, err error) {
	s, err := newScanner[M](rows)
	if err != nil {
		return nil, err
	}

	for s.rows.Next() {
		msg, err := s.scanRow()
		if err != nil {
			return nil, err
		}

		result = append(result, msg)
	}

	return result, err
}

// ScanOne returns a single instance of proto Message with type M, filled with data from rows.
// pgx.ErrNoRows is returned when there are rows to scan.
// See Scan for field name matching rules.
func ScanOne[M proto.Message](rows pgx.Rows) (M, error) {
	var m M

	s, err := newScanner[M](rows)
	if err != nil {
		return m, err
	}

	if !s.rows.Next() {
		return m, pgx.ErrNoRows
	}

	return s.scanRow()
}

type ServerStream[M proto.Message] interface {
	Send(M) error
	Context() context.Context
}

// ScanStream writes instances of proto messages with type M to stream.Send(), filled with data from rows.
// ScanStream returns a nil error when rows is exhausted or an error when one is encountered,
// durng scanning or sending.
// Messages may already have been send when returning an error.
// See Scan for field name matching rules.
func ScanStream[M proto.Message](rows pgx.Rows, stream ServerStream[M]) error {
	s, err := newScanner[M](rows)
	if err != nil {
		return err
	}

	for s.rows.Next() {
		msg, err := s.scanRow()
		if err != nil {
			return err
		}

		if err := stream.Send(msg); err != nil {
			return fmt.Errorf("pbpgx.ScanStream: %w", err)
		}
	}

	return nil
}
