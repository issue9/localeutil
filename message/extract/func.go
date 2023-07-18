// SPDX-License-Identifier: MIT

package extract

import (
	"go/ast"
	"path"
	"strings"
)

// 表示本地化的函数
type localeFunc struct {
	path string // 函数的完整导入路径
	name string // 函数名
}

// 表示由 import 转换后的函数名
type importFunc struct {
	modName string // import 中的别名
	name    string // 函数名
}

func split(funcs ...string) []localeFunc {
	ret := make([]localeFunc, 0, len(funcs))
	for _, f := range funcs {
		if index := strings.LastIndexByte(f, '.'); index > 0 {
			ret = append(ret, localeFunc{path: f[:index], name: f[index+1:]})
		}
	}
	return ret
}

func filterImportFuncs(imports []*ast.ImportSpec, funcs []localeFunc) []importFunc {
	mods := make([]importFunc, 0, len(funcs))

	for _, ip := range imports {
		var modName string
		if ip.Name != nil && ip.Name.Name != "" {
			modName = ip.Name.Name
		} else {
			modName = path.Base(strings.Trim(ip.Path.Value, "\""))
		}

		modPath := strings.Trim(ip.Path.Value, "\"")

		for _, f := range funcs {
			if f.path != modPath {
				continue
			}

			mods = append(mods, importFunc{modName: modName, name: f.name})

		}
	}

	return mods
}
