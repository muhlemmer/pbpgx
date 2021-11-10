// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

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
