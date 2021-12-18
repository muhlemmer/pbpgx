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
	"testing"

	"github.com/muhlemmer/pbpgx/internal/support"
	"github.com/muhlemmer/pbpgx/internal/testlib"
	"google.golang.org/protobuf/proto"
)

func TestReadOne(t *testing.T) {
	tests := []struct {
		name  string
		query *support.SimpleQuery
		want  *support.Simple
	}{
		{
			"nil columns",
			&support.SimpleQuery{
				Id: 5,
			},
			&support.Simple{
				Id:    5,
				Title: "five",
				Data:  "five is a four letter word",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tab := NewTable("public", "simple_ro", nil)

			got, err := ReadOne[*support.Simple, int32, support.SimpleColumns](testlib.CTX, testlib.ConnPool, tab, tt.query)
			if err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("ReadOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
