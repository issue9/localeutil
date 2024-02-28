// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

//go:build windows

// Package windows 与 windows 系统相关的一些本地化操作函数
package windows

import (
	"syscall"
	"unsafe"

	"github.com/issue9/localeutil/internal/dll"
)

// GetLCID 将一个符合 BCP47 的名称转换成 LCID
func GetLCID(bcp string) (uint32, error) {
	f := dll.Kernel32().NewProc("LocaleNameToLCID")
	name, err := syscall.UTF16FromString(bcp)
	if err != nil {
		return 0, err
	}

	r1, _, err := f.Call(uintptr(unsafe.Pointer(&name[0])), uintptr(0))
	if r1 != 0 {
		return uint32(r1), nil
	}
	return 0, err
}
