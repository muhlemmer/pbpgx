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
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
)

func Test_bufLenCache(t *testing.T) {
	var c bufLenCache

	tests := []struct {
		name string
		set  int
		want int
	}{
		{
			"One",
			1,
			1,
		},
		{
			"Higher",
			10,
			10,
		},
		{
			"Lower",
			9,
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.setHigher(tt.set)
			if got := c.get(); got != tt.want {
				t.Errorf("bufLenCache.get() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestOnEmpty_pgStatus(t *testing.T) {
	tests := []struct {
		name string
		o    OnEmpty
		want pgtype.Status
	}{
		{
			"zero",
			Zero,
			pgtype.Present,
		},
		{
			"null",
			Null,
			pgtype.Null,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.pgStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnEmpty.pgStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
