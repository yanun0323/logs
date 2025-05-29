// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package buffer provides a pool-allocated byte buffer.
package buffer

import (
	"bytes"
	"sync"
)

const (
	// 初始 buffer 大小，根據典型日誌長度調整
	defaultBufferSize = 256
	// 最大 buffer 大小，避免記憶體浪費
	maxBufferSize = 64 * 1024 // 64KB
)

// Having an initial size gives a dramatic speedup.
var Pool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, defaultBufferSize)
		return bytes.NewBuffer(buf)
	},
}

// Put 安全地將 buffer 放回 pool
func Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}

	// 如果 buffer 太大，不要放回 pool 以節省記憶體
	if buf.Cap() > maxBufferSize {
		return
	}

	buf.Reset()
	Pool.Put(buf)
}

// Get 從 pool 獲取 buffer
func Get() *bytes.Buffer {
	return Pool.Get().(*bytes.Buffer)
}
