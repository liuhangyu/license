// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd openbsd netbsd dragonfly

package fsnotify

import 	"code.uni-ledger.com/switch/license/public/deplib/golang.org/x/sys/unix"
const openMode = unix.O_NONBLOCK | unix.O_RDONLY
