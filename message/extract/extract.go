// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package extract 提供从 Go 源码中提取本地化内容的功能
package extract

import (
	"cmp"
	"context"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/issue9/localeutil"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/packages"

	"github.com/issue9/localeutil/message"
)

const mode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
	packages.NeedImports | packages.NeedModule | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo

type extractor struct {
	warnLog message.LogFunc
	infoLog message.LogFunc
	fset    *token.FileSet
	funcs   []fn
	root    string

	mux sync.Mutex
	msg []message.Message
}

// Extract 提取本地化内容
//
// o 给定的参数错误，可能会触发 panic，比如 o 为空、o.Funcs 格式错误等。
func Extract(ctx context.Context, o *Options) (*message.File, error) {
	if o == nil {
		panic("参数 o 不能为空")
	}

	ex, err := o.buildExtractor()
	if err != nil {
		return nil, err
	}

	dirs, err := getDir(ex.root, o.Recursive, o.SkipSubModule)
	if err != nil {
		return nil, err
	}

	if err := ex.inspectDirs(ctx, dirs); err != nil {
		return nil, err
	}

	slices.SortStableFunc(ex.msg, func(a, b message.Message) int { return cmp.Compare(a.Key, b.Key) })

	return &message.File{Languages: []language.Tag{o.Language}, Messages: ex.msg}, nil
}

func (ex *extractor) inspectDirs(ctx context.Context, dirs []string) error {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for _, dir := range dirs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wg.Add(1)
			go func(dir string) {
				defer wg.Done()
				ex.inspectDir(ctx, dir)
			}(dir)
		}
	}

	return nil
}

func (ex *extractor) inspectDir(ctx context.Context, dir string) error {
	cfg := &packages.Config{
		Mode:    mode,
		Context: ctx,
		Dir:     dir,
		Fset:    ex.fset,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	defer wg.Wait()
	for _, pkg := range pkgs {
		info := pkg.TypesInfo
		for _, f := range pkg.Syntax {
			wg.Add(1)
			go func(f *ast.File) {
				defer wg.Done()

				ast.Inspect(f, func(n ast.Node) bool {
					switch expr := n.(type) {
					case *ast.TypeSpec, *ast.ImportSpec:
						return false
					case *ast.CallExpr:
						return ex.inspect(expr, info)
					default:
						return true
					}
				})
			}(f)
		}
	}

	return nil
}

// 遍历 expr 表达式
//
// 返回值表示是否需要访问子元素
func (ex *extractor) inspect(expr *ast.CallExpr, info *types.Info) bool {
	t := info.TypeOf(expr.Fun)
	switch typ := t.(type) {
	case *types.Signature: // 所有 () 形式的调用
		if typ.Params().Len() == 0 || typ.Params().At(0).Type() != types.Typ[types.String] { // 可能是匿名函数
			return true
		}

		var obj types.Object
		switch ft := expr.Fun.(type) {
		case *ast.SelectorExpr:
			obj = info.ObjectOf(ft.Sel)
		case *ast.Ident:
			obj = info.ObjectOf(ft)
		}
		if obj == nil {
			return true
		}
		f, ok := obj.(*types.Func)
		if !ok {
			return true
		}

		s := f.Signature()   // typ.Recv 永远返回 nil，只有通过 types.Func.Signature 返回的才会有正确的返回值
		if s.Recv() == nil { // func
			if !ex.tryAppendMsg(expr, f.Pkg().Path(), "", f.Name()) {
				return false
			}
		} else { // method
			pkgName, structName := parseTypeName(s.Recv().Type().String())
			if !ex.tryAppendMsg(expr, pkgName, structName, f.Name()) {
				return false
			}
		}
	case *types.Alias: // type Alias = localeutil.StringPhrase; Alias('key')
		rhs := typ.Rhs()
		alias, ok := rhs.(*types.Alias)
		for ok {
			rhs = alias.Rhs()
			if _, bok := rhs.(*types.Basic); bok {
				rhs = alias
				break
			}
			alias, ok = rhs.(*types.Alias)
		}

		pkgName, funcName := parseTypeName(rhs.String())
		if !ex.tryAppendMsg(expr, pkgName, "", funcName) {
			return false
		}
	case *types.Named: // type X string; X('key')
		obj := typ.Obj()
		if !ex.tryAppendMsg(expr, obj.Pkg().Path(), "", obj.Name()) {
			return false
		}
	case *types.Basic:
		return false
	}
	return true
}

func (ex *extractor) tryAppendMsg(expr *ast.CallExpr, pkgName, structName, name string) (continueInspect bool) {
	index := slices.IndexFunc(ex.funcs, func(m fn) bool {
		return m.name == name && pkgName == m.pkgName && structName == m.typeName
	})
	if index < 0 {
		return true
	}

	ex.appendMsg(expr)
	return false
}

func (ex *extractor) appendMsg(expr *ast.CallExpr) {
	var key string
	p := ex.fset.Position(expr.Pos())
	path := ex.trimPath(p.Filename)

	switch v := expr.Args[0].(type) {
	case *ast.BasicLit: // 直接参数，比如 call("xxx")
		key = v.Value
	case *ast.Ident: // 间接参数，比如：const xxx; call(xxx) 或是 var xxx; call(xxx)
		switch d := v.Obj.Decl.(type) {
		case *ast.ValueSpec:
			if d.Names != nil && d.Names[0].Obj.Kind == ast.Con { // 常量，可获得值
				key = d.Values[0].(*ast.BasicLit).Value
			} else { // 变量，编译时无法获得
				pos := ex.fset.Position(expr.Pos())
				file := ex.trimPath(pos.Filename)
				ex.warnLog(localeutil.Phrase("can not covert to message at %s:%d", file, pos.Line))
				return
			}
		}
	default:
		ex.warnLog(localeutil.Phrase("can not covert to message at %s:%d", path, p.Line))
		return
	}

	if key != "" {
		key = key[1 : len(key)-1]
	}
	if key == "" {
		ex.warnLog(localeutil.Phrase("has empty string at %s:%d", path, p.Line))
		return
	}

	ex.mux.Lock()
	defer ex.mux.Unlock()

	if slices.IndexFunc(ex.msg, func(m message.Message) bool { return m.Key == key }) >= 0 {
		ex.warnLog(localeutil.Phrase("has same key %s at %s:%d, will be ignore", strconv.Quote(key), path, p.Line))
		return
	}

	ex.infoLog(localeutil.Phrase("find new locale string %s at %s:%d", strconv.Quote(key), path, p.Line))
	ex.msg = append(ex.msg, message.Message{Key: key, Message: message.Text{Msg: key}})
}

func parseTypeName(t string) (pkg, structure string) {
	if t[0] == '*' {
		t = t[1:]
	}

	if index := strings.LastIndexByte(t, '['); index >= 0 {
		t = t[:index]
	}

	if index := strings.LastIndexByte(t, '.'); index >= 0 {
		pkg = t[:index]
		t = t[index+1:]
	}

	return pkg, t
}

func (ex *extractor) trimPath(p string) string {
	path := strings.TrimPrefix(p, ex.root) // 只显示相对于检测目录的路径
	if path != "" && (path[0] == '/' || path[0] == filepath.Separator) {
		path = path[1:]
	}
	if path == "" {
		return "./"
	}
	return path
}
