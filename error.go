// Copyright (c) 2021, Tim Möhlmann. All rights reserved.
// Use of this source code is governed by a License that can be found in the LICENSE file.
// SPDX-License-Identifier: BSD-3-Clause

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
