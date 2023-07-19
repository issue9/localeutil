// SPDX-License-Identifier: MIT

package extract

import (
	"fmt"
	"go/ast"
	"path"
	"strings"
)

// 表示本地化的函数
type localeFunc struct {
	path      string // 函数的完整导入路径
	structure string // 类型名，可能为空
	name      string // 函数名
}

// 表示由 import 转换后的函数名
type importFunc struct {
	modName    string // import 中的别名
	structName string // 类型名，可能为空
	name       string // 函数名
}

func split(funcs ...string) []localeFunc {
	ret := make([]localeFunc, 0, len(funcs))
	for _, f := range funcs {
		base := path.Base(f)
		dir := path.Dir(f)
		switch strs := strings.Split(base, "."); len(strs) {
		case 2:
			ret = append(ret, localeFunc{path: path.Join(dir, strs[0]), name: strs[1]})
		case 3:
			ret = append(ret, localeFunc{path: path.Join(dir, strs[0]), structure: strs[1], name: strs[2]})
		default:
			panic(fmt.Sprintf("%s 格式无效", f))
		}
	}
	return ret
}

func filterImportFuncs(fileModPath string, imports []*ast.ImportSpec, funcs []localeFunc) []importFunc {
	mods := make([]importFunc, 0, len(funcs))

	for _, f := range funcs {
		if fileModPath == f.path {
			mods = append(mods, importFunc{name: f.name, structName: f.structure})
			continue
		}

		for _, ip := range imports {
			modPath := strings.Trim(ip.Path.Value, "\"")
			if f.path != modPath {
				continue
			}

			var modName string
			if ip.Name != nil && ip.Name.Name != "" {
				modName = ip.Name.Name
			} else {
				modName = path.Base(modPath)
			}

			mods = append(mods, importFunc{modName: modName, name: f.name, structName: f.structure})
		}
	}

	return mods
}
