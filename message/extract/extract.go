// SPDX-License-Identifier: MIT

// Package extract 提供从 Go 源码中提取本地化内容的功能
package extract

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/issue9/sliceutil"
	"golang.org/x/text/language"

	"github.com/issue9/localeutil/message"
)

// Logger 日志输出接口
type Logger interface {
	Print(...any)
	Printf(string, ...any)
}

type extracter struct {
	log   Logger
	funcs []localeFunc
	fset  *token.FileSet

	mux sync.Mutex
	msg []message.Message
}

// Extract 提取本地化内容
//
// lang 代码的文本所使用的语言；
// root 需要提取本地化内容的源码目录；
// f 表示被用于本地化的函数，所有 f 指定的函数，其参数都将被提取为本地化的内容。
// f 每个元素的格式为 mod/path.func，mod/path 为包的导出路径，func 为函数名。
// f 的函数签名应该始终与 [localeutil.Phrase] 相同；
func Extract(ctx context.Context, lang, root string, r bool, log Logger, f ...string) (*message.Messages, error) {
	// NOTE: 有可能存在将 localeutil.Phrase 二次封装的情况，
	// 为了尽可能多地找到本地化字符串，所以采用用户指定函数的方法。

	dirs, err := getDir(root, r)
	if err != nil {
		return nil, err
	}

	ex := &extracter{
		log:   log,
		funcs: split(f...),
		fset:  token.NewFileSet(),

		msg: make([]message.Message, 0, 100),
	}

	if err := ex.scanDirs(ctx, dirs); err != nil {
		return nil, err
	}

	sort.SliceStable(ex.msg, func(i, j int) bool { return ex.msg[i].Key < ex.msg[j].Key })

	return &message.Messages{
		Languages: []*message.Language{
			{ID: language.MustParse(lang), Messages: ex.msg},
		},
	}, nil
}

func (ex *extracter) scanDirs(ctx context.Context, dirs []string) error {
	wg := &sync.WaitGroup{}
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, e := range entries {
			select {
			case <-ctx.Done():
				return context.Canceled
			default:
				name := strings.ToLower(e.Name())
				if e.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
					continue
				}

				wg.Add(1)
				go func(p string) {
					defer wg.Done()

					f, err := parser.ParseFile(ex.fset, p, nil, parser.ParseComments)
					if err != nil {
						ex.log.Print(err)
						return
					}

					ex.inspectFile(f)
				}(filepath.Join(dir, e.Name()))
			}
		}
	}
	wg.Wait()

	return nil
}

func (ex *extracter) inspectFile(f *ast.File) {
	mods := filterImportFuncs(f.Imports, ex.funcs)
	ast.Inspect(f, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.TypeSpec, *ast.ImportSpec:
			return false
		case *ast.CallExpr:
			msg := inspect(ex.fset, expr, mods, ex.log)
			if msg.Key == "" {
				return true
			}

			ex.mux.Lock()
			defer ex.mux.Unlock()

			if sliceutil.Exists(ex.msg, func(m message.Message) bool { return m.Key == msg.Key }) {
				p := ex.fset.Position(expr.Pos())
				log.Printf("存在相同的本地化信息 %s，将被忽略，位于：%s:%d", msg.Key, p.Filename, p.Line)
				return true
			}
			ex.msg = append(ex.msg, msg)
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

	var key string
	switch v := expr.Args[0].(type) {
	case *ast.BasicLit:
		key = v.Value
	case *ast.Ident: // 常量/变量
		switch d := v.Obj.Decl.(type) {
		case *ast.ValueSpec:
			if d.Names != nil && d.Names[0].Obj.Kind == ast.Con {
				key = d.Values[0].(*ast.BasicLit).Value
			} else {
				log.Printf("当前类型 %s 无法转换成本地化信息", d.Names[0].Obj.Kind)
			}
		}
	}

	if key == "" {
		return msg
	}
	key = key[1 : len(key)-1]
	msg.Key = key
	msg.Message.Msg = key
	return msg
}
