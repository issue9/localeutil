// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package syslocale

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/issue9/localeutil/internal/dll"
)

func getOSLocaleName() string {
	f := dll.Kernel32().NewProc("GetUserDefaultLocaleName")

	const maxLen = 85 // GetUserDefaultLocaleName 第二个参数
	buf := make([]uint16, maxLen)
	r1, _, err := f.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(maxLen))
	if uint32(r1) == 0 {
		log.Println(err)
		return ""
	}
	return syscall.UTF16ToString(buf)
}
