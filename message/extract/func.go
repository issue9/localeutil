// SPDX-License-Identifier: MIT

package extract

import (
	"fmt"
	"go/ast"
	"path"
	"strings"
)

// 表示由 import 转换后的函数名
type importFunc struct {
	modName    string // import 中的别名
	structName string // 类型名，可能为空
	name       string // 函数名
}

func filterImportFuncs(fileModPath string, imports []*ast.ImportSpec, funcList []string) []importFunc {
	funcs := split(funcList...)

	mods := make([]importFunc, 0, len(funcs))

	for _, f := range funcs {
		if fileModPath == f.modName {
			mods = append(mods, importFunc{name: f.name, structName: f.structName})
			continue
		}

		for _, ip := range imports {
			modPath := strings.Trim(ip.Path.Value, "\"")
			if f.modName != modPath {
				continue
			}

			var modName string
			if ip.Name != nil && ip.Name.Name != "" {
				modName = ip.Name.Name
			} else {
				modName = path.Base(modPath)
			}

			mods = append(mods, importFunc{modName: modName, name: f.name, structName: f.structName})
		}
	}

	return mods
}

// 返回从 [Options.Funcs] 中分析而来的中间数据
// 此时返回元素中的 modName 表示的是完整的模块导出山路径。
func split(funcs ...string) []importFunc {
	ret := make([]importFunc, 0, len(funcs))
	for _, f := range funcs {
		base := path.Base(f)
		dir := path.Dir(f)
		switch strs := strings.Split(base, "."); len(strs) {
		case 2:
			ret = append(ret, importFunc{modName: path.Join(dir, strs[0]), name: strs[1]})
		case 3:
			ret = append(ret, importFunc{modName: path.Join(dir, strs[0]), structName: strs[1], name: strs[2]})
		default:
			panic(fmt.Sprintf("%s 格式无效", f))
		}
	}
	return ret
}
