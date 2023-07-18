// SPDX-License-Identifier: MIT

package extract

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/issue9/localeutil/message"
	"github.com/issue9/sliceutil"
	"golang.org/x/text/language"
)

// Logger 日志输出接口
type Logger interface {
	Printf(string, ...any)
}

// Extract 提取本地化内容
//
// lang 代码的文本所使用的语言；
// root 需要提取本地化内容的源码目录；
// f 表示被用于本地化的函数，所有 f 指定的函数，其参数都将被提取为本地化的内容。
// f 每个元素的格式为 mod/path.func，mod/path 为包的导出路径，func 为函数名；
//
// 如果源码错误，返回该错误信息。如果是因为不能序列化，通过 log 输出错误信息；
func Extract(ctx context.Context, lang, root string, r bool, log Logger, f ...string) (*message.Messages, error) {
	// NOTE: 有可能存在将 localeutil.Phrase 二次封装的情况，
	// 为了尽可能多地找到本地化字符串，所以采用用户指定函数的方法。

	dirs, err := getDir(root, r)
	if err != nil {
		return nil, err
	}

	funs := split(f...)
	fset := token.NewFileSet()
	l := &message.Language{ID: language.MustParse(lang), Messages: make([]message.Message, 0, 100)}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if e.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
				continue
			}

			p := filepath.Join(dir, e.Name())
			f, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
			if err != nil {
				return nil, err
			}

			inspectFile(fset, f, funs, l, log)
		}
	}

	return &message.Messages{
		Languages: []*message.Language{l},
	}, nil
}

func inspectFile(fset *token.FileSet, f *ast.File, funs []localeFunc, l *message.Language, log Logger) {
	mods := filterImportFuncs(f.Imports, funs)
	ast.Inspect(f, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.TypeSpec, *ast.ImportSpec:
			return false
		case *ast.CallExpr:
			msg := inspect(fset, expr, mods, log)
			if msg.Key == "" {
				return true
			}

			if sliceutil.Exists(l.Messages, func(m message.Message) bool { return m.Key == msg.Key }) {
				p := fset.Position(expr.Pos())
				log.Printf("存在相同的本地化信息 %s，将被忽略，位于：%s:%d", msg.Key, p.Filename, p.Line)
				return true
			}
			l.Messages = append(l.Messages, msg)
			return false
		}
		return true
	})
}

func inspect(fset *token.FileSet, expr *ast.CallExpr, mods []importFunc, log Logger) message.Message {
	msg := message.Message{}

	f, ok := expr.Fun.(*ast.SelectorExpr)
	if !ok {
		return msg
	}

	var xName string
	switch ft := f.X.(type) {
	case *ast.CallExpr: // localeutil.Phrase(xxx).LocaleString(p)
		return inspect(fset, ft, mods, log)
	case *ast.Ident:
		xName = ft.Name
	default:
		return msg
	}

	exists := sliceutil.Exists(mods, func(m importFunc) bool {
		return xName == m.modName && m.name == f.Sel.Name
	})
	if !exists {
		return msg
	}

	switch v := expr.Args[0].(type) {
	case *ast.BasicLit:
		msg.Key = v.Value
		msg.Message.Msg = v.Value
	case *ast.Ident: // 常量/变量
		switch d := v.Obj.Decl.(type) {
		case *ast.ValueSpec:
			if d.Names != nil && d.Names[0].Obj.Kind == ast.Con {
				msg.Key = d.Values[0].(*ast.BasicLit).Value
				msg.Message.Msg = msg.Key
			} else {
				log.Printf("当前类型 %s 无法转换成本地化信息", d.Names[0].Obj.Kind)
			}
		}
	}

	return msg
}
