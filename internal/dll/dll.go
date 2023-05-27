// SPDX-License-Identifier: MIT

// Package dll windows 平台下一些 dll
package dll

import "syscall"

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

func Kernel32() *syscall.LazyDLL { return kernel32 }
