// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && (arm || arm64)

package mobileinit

import (
	"log"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <asl.h>
#include <stdlib.h>

void log_wrap(const char *logStr);
*/
import "C"

type aslWriter struct{}

func (aslWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.log_wrap(cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}

func init() {
	log.SetOutput(aslWriter{})
}
