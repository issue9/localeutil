// SPDX-License-Identifier: MIT

package testdata

import (
	"strings"

	"github.com/issue9/localeutil"
	l "github.com/issue9/localeutil"
)

const err = "error "

var (
	_ = l.Phrase("alias 1")
	_ = localeutil.Error(err)
	_ = localeutil.Error(err + "1")              // 不支持计算
	_ = localeutil.Error(strings.TrimSpace(err)) // 无法获取
)
