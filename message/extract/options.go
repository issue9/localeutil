// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package extract

import (
	"fmt"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"

	"github.com/issue9/localeutil"
	"github.com/issue9/localeutil/message"
)

type Options struct {
	// Language 提取内容的语言 ID
	Language language.Tag

	// 读取的根目录
	//
	// 需要位于一个 Go 的模块中。
	Root string

	// 是否读取子目录的内容
	Recursive bool

	// 忽略子模块
	//
	// 当 Recursive 为 true 时，此值为 true，可以不读取子模块的内容。
	SkipSubModule bool

	// 警告日志通道
	//
	// 默认为输出到终端。
	WarnLog message.LogFunc

	// 普通信息日志通道
	//
	// 主要报告提取的进度，如果为空，则不输出内容。
	InfoLog message.LogFunc

	// 用于提取本地化内容的函数列表
	//
	// 每个元素的格式为：
	//  mod/path[.type].func
	//
	// mod/path 为包的导出路径；
	//
	// type 为类型名称，可以省略；
	//
	// func 为用于实现本地化的调用，可能是与 type 关联的方法
	// 或是无 type 的函数还有可能是简单的类型转换。
	// func 至少需要一个参数，且其第一个参数的类型必须为 string。
	//
	// 能正确识别别名，比如：
	//  type x = localeutil.Printer
	// 当在 Funcs 指定了 github.com/issue9/localeutil.Printer 时，也会识别 x。
	//
	// 如果指定的是接口类型的方法，在提取时不会主动判断是否实现了该接口，
	// 必须在代码中明确为该接口类型的方法才会被提取：
	//  var p PrinterInterface = &Printer{}
	//  p.Printf(...) // 正常提取内容
	//  var p = &Printer{}
	//  p.Printf(...) // 即使 Printer 实现了 PrinterInterface 接口，也不会提取内容。
	//
	// 如果为空将不输出任何内容，格式错误将会触发 panic。
	Funcs []string
}

// [Options.Funcs] 转换后的表示
type fn struct {
	pkgName  string // 包名
	typeName string // 类型名，可能为空
	name     string // 函数名
}

func (o *Options) buildExtractor() (*extractor, error) {
	if o.WarnLog == nil {
		o.WarnLog = func(v localeutil.Stringer) { log.Println(v) } // TODO(go1.21): slog
	}
	if o.InfoLog == nil {
		o.InfoLog = func(v localeutil.Stringer) { log.Println(v) }
	}

	abs, err := filepath.Abs(o.Root)
	if err != nil {
		return nil, err
	}
	o.Root = abs

	return &extractor{
		warnLog: o.WarnLog,
		infoLog: o.InfoLog,
		fset:    token.NewFileSet(),
		funcs:   split(o.Funcs...),
		root:    abs,

		msg: make([]message.Message, 0, 100),
	}, nil
}

func getDir(root string, r, skip bool) ([]string, error) {
	if !r {
		return []string{root}, nil
	}

	dirs := make([]string, 0, 30)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		if skip && root != p {
			stat, err := os.Stat(filepath.Join(p, "go.mod"))
			if err == nil && !stat.IsDir() {
				return fs.SkipDir
			}
		}

		dirs = append(dirs, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

// 返回从 [Options.Funcs] 中分析而来的中间数据
// 此时返回元素中的 pkgName 表示的是完整的模块导出山路径。
func split(funcs ...string) []fn {
	ret := make([]fn, 0, len(funcs))
	for _, f := range funcs {
		base := path.Base(f)
		dir := path.Dir(f)
		switch strs := strings.Split(base, "."); len(strs) {
		case 2:
			ret = append(ret, fn{pkgName: path.Join(dir, strs[0]), name: strs[1]})
		case 3:
			ret = append(ret, fn{pkgName: path.Join(dir, strs[0]), typeName: strs[1], name: strs[2]})
		default:
			panic(fmt.Sprintf("%s 格式无效", f))
		}
	}
	return ret
}
