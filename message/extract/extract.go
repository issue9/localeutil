// SPDX-License-Identifier: MIT

// Package extract 提供从 Go 源码中提取本地化内容的功能
package extract

import (
	"context"
	"errors"
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
	"github.com/issue9/source"
	"golang.org/x/text/language"

	"github.com/issue9/localeutil/message"
)

// Logger 日志输出接口
type Logger interface {
	Print(...any)
	Printf(string, ...any)
}

type Options struct {
	// Language 提取内容的语言 ID
	Language string

	// 读取的根目录
	//
	// 需要位于一个 Go 的模块中。
	Root string

	// 是否读取子目录的内容
	Recursive bool

	// 日志输出通道
	Log Logger

	// 用于输出本地化内容的函数列表
	//
	// 每个元素的格式为 mod/path[.struct].func，mod/path 为包的导出路径，
	// struct 为结构体名称，可以省略，func 为函数或方法名。
	//
	// 函数至少需要一个参数，且其第一个参数的类型必须为 string。
	// 如果指向的是方法，那么在调用此方法的结构必须有明确类型声明，不能由类型推荐获得。
	// 比如，当 f 为 golang.org/x/text/message.Printer.Printf 时：
	//
	//	// 以下无法提取内容
	//	p := message.NewPrinter();
	//	p.Printf(...)
	//
	//	// 以下可以
	//	var p *message.Printer = message.NewPrinter();
	//	p.Printf(...)
	Funcs []string
}

type extracter struct {
	log   Logger
	funcs []localeFunc
	fset  *token.FileSet

	mux sync.Mutex
	msg []message.Message
}

// Extract 提取本地化内容
func Extract(ctx context.Context, o *Options) (*message.Messages, error) {
	// NOTE: 有可能存在将 localeutil.Phrase 二次封装的情况，
	// 为了尽可能多地找到本地化字符串，所以采用用户指定函数的方法。

	dirs, err := getDir(o.Root, o.Recursive)
	if err != nil {
		return nil, err
	}

	ex := &extracter{
		log:   o.Log,
		funcs: split(o.Funcs...),
		fset:  token.NewFileSet(),

		msg: make([]message.Message, 0, 100),
	}

	if err := ex.scanDirs(ctx, dirs); err != nil {
		return nil, err
	}

	sort.SliceStable(ex.msg, func(i, j int) bool { return ex.msg[i].Key < ex.msg[j].Key })

	return &message.Messages{
		Languages: []*message.Language{
			{ID: language.MustParse(o.Language), Messages: ex.msg},
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

					ex.inspectFile(p, f)
				}(filepath.Join(dir, e.Name()))
			}
		}
	}
	wg.Wait()

	return nil
}

func (ex *extracter) inspectFile(p string, f *ast.File) {
	modPath, err := source.ModPath(p)
	switch {
	case errors.Is(err, os.ErrNotExist):
		ex.log.Printf("必须是在一个模块中，当前无法找到 go.mod 文件")
		return
	case err != nil:
		ex.log.Print(err)
		return
	}

	mods := filterImportFuncs(modPath, f.Imports, ex.funcs)
	ast.Inspect(f, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.TypeSpec, *ast.ImportSpec:
			return false
		case *ast.CallExpr:
			msg := ex.inspect(expr, mods)
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

func (ex *extracter) inspect(expr *ast.CallExpr, mods []importFunc) message.Message {
	msg := message.Message{}
	var modName, structName, name string

	switch f := expr.Fun.(type) {
	case *ast.SelectorExpr:
		switch ft := f.X.(type) {
		case *ast.CallExpr: // localeutil.Phrase(xxx).LocaleString(p)
			return ex.inspect(ft, mods)
		case *ast.Ident:
			if ft.Obj != nil {
				modName, structName = ex.getObjectName(ft.Obj)
			} else {
				modName = ft.Name
			}
			name = f.Sel.Name
		default:
			return msg
		}
	case *ast.Ident:
		name = f.Name
	}

	exists := sliceutil.Exists(mods, func(m importFunc) bool {
		ok := m.name == name && modName == m.modName
		if structName != "" {
			ok = ok && structName == m.structName
		}

		return ok
	})
	if !exists {
		return msg
	}

	var key string
	switch v := expr.Args[0].(type) {
	case *ast.BasicLit:
		key = v.Value
	case *ast.Ident: // const / var
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

func (ex *extracter) getObjectName(obj *ast.Object) (modName, structName string) {
	switch decl := obj.Decl.(type) {
	case *ast.ValueSpec: // 局部变量/全局变量
		if decl.Type != nil {
			return getExprNames(decl.Type, ex.log)
		}
	case *ast.Field: // 函数参数
		return getExprNames(decl.Type, ex.log)
	}
	return "", ""
}

func getExprNames(expr ast.Expr, log Logger) (modName, structName string) {
	switch s := expr.(type) {
	case *ast.SelectorExpr:
		return s.X.(*ast.Ident).Name, s.Sel.Name
	case *ast.Ident:
		return "", s.Name
	case *ast.StarExpr:
		return getExprNames(s.X, log)
	}
	return "", ""
}
