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
	"runtime"
)

// checkNotRuntimeErr checks the return value from recover().
// If the type is a runtime.Error of unknown, a call to panic is made,
// to continue any runtime panic.
// In case the pass value is any other error type, it is returned.
func checkNotRuntimeErr(v interface{}) error {
	switch err := v.(type) {
	case nil:
		return nil
	case runtime.Error:
		panic(err)
	case error:
		return err
	default:
		panic(err)
	}
}
