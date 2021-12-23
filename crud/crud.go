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

// Package crud provides Create Read Update and Delete functionality on top of pbpgx.
// It is capable of effeciently building queries based on incomming protocol buffer messages,
// and returning results as protocol buffer messages, using protoreflect and generics.
package crud

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/muhlemmer/pbpgx"
)

type Preparator interface {
	pbpgx.Executor
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}
