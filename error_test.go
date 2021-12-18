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
	"errors"
	"testing"
)

func Test_checkNotRuntimeErr(t *testing.T) {
	tests := []struct {
		name      string
		f         func()
		wantErr   bool
		wantPanic bool
	}{
		{
			"nil",
			func() {},
			false,
			false,
		},
		{
			"runtime Error",
			func() {
				var i []interface{}
				_ = i[0]
			},
			false,
			true,
		},
		{
			"some error",
			func() { panic(errors.New("foo")) },
			true,
			false,
		},
		{
			"Panic with non-error type",
			func() { panic("foo") },
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			if tt.wantPanic {
				defer func() {
					if v := recover(); v == nil {
						t.Error("Expected panic")
					}
				}()
			}

			func() {
				defer func() { err = checkNotRuntimeErr(recover()) }()
				tt.f()
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("checkNotRuntimeErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}
}
