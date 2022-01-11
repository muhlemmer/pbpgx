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
	"testing"

	"github.com/jackc/pgtype"
	"github.com/muhlemmer/pbpgx/internal/support"
)

// Lazy test for coverage and nil error only.
// Correct functioning is confirmed by pbpgx.Scan tests.
func TestNew(t *testing.T) {
	fields := new(support.Supported).ProtoReflect().Descriptor().Fields()
	n := fields.Len()

	dest := make([]interface{}, n)

	for i := 0; i < n; i++ {
		d, err := New(fields.Get(i), pgtype.Undefined)
		if err != nil {
			t.Fatal(err)
		}
		dest[i] = d
	}
}

func TestNew_error(t *testing.T) {
	fields := new(support.Unsupported).ProtoReflect().Descriptor().Fields()
	n := fields.Len()

	dest := make([]interface{}, n)

	for i := 0; i < n; i++ {
		d, err := New(fields.Get(i), pgtype.Undefined)
		if err == nil {
			t.Fatal("New: expected error, got nil")
		}
		dest[i] = d
	}
}
