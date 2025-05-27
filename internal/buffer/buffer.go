// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package buffer provides a pool-allocated byte buffer.
package buffer

import (
	"bytes"
	"sync"
)

// Having an initial size gives a dramatic speedup.
var Pool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}
