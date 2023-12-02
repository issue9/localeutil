// SPDX-License-Identifier: MIT

package extract

import (
	"fmt"
	"go/token"
	"io/fs"
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
	Log message.LogFunc

	// 用于输出本地化内容的函数列表
	//
	// 每个元素的格式为 mod/path[.struct].func，mod/path 为包的导出路径，
	// struct 为结构体名称，可以省略，func 为函数或方法名或是类型。
	//
	// func 至少需要一个参数，且其第一个参数的类型必须为 string。
	// 如果指定的是接口类型的方法，必须明确类型与接口相同的对象才会提取内容：
	//  var p PrinterInterface = &Printer{}
	//  p.Printf(...) // 正常提取内容
	//  var p = &Printer{}
	//  p.Printf(...) // 即使 Printer 实现了 PrinterInterface 接口，也不会提取内容。
	//
	// 如果为空将不输出任何内容，格式错误将会触发 panic。
	Funcs []string
}

// [Options.Funcs] 转换后的表示
type importFunc struct {
	pkgName    string // 包名
	structName string // 类型名，可能为空
	name       string // 函数名
}

func (o *Options) buildExtractor() (*extractor, error) {
	if o.Log == nil {
		o.Log = func(v localeutil.Stringer) { fmt.Print(v) }
	}

	return &extractor{
		log:   o.Log,
		fset:  token.NewFileSet(),
		funcs: split(o.Funcs...),

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
// 此时返回元素中的 modName 表示的是完整的模块导出山路径。
func split(funcs ...string) []importFunc {
	ret := make([]importFunc, 0, len(funcs))
	for _, f := range funcs {
		base := path.Base(f)
		dir := path.Dir(f)
		switch strs := strings.Split(base, "."); len(strs) {
		case 2:
			ret = append(ret, importFunc{pkgName: path.Join(dir, strs[0]), name: strs[1]})
		case 3:
			ret = append(ret, importFunc{pkgName: path.Join(dir, strs[0]), structName: strs[1], name: strs[2]})
		default:
			panic(fmt.Sprintf("%s 格式无效", f))
		}
	}
	return ret
}
