// SPDX-License-Identifier: MIT

// Package extract 提供从 Go 源码中提取本地化内容的功能
package extract

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"
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
func Extract(ctx context.Context, o *Options) (*message.Language, error) {
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

	sort.SliceStable(ex.msg, func(i, j int) bool { return ex.msg[i].Key < ex.msg[j].Key })

	return &message.Language{ID: o.Language, Messages: ex.msg}, nil
}

func (ex *extractor) info(msg localeutil.Stringer) {
	if ex.infoLog != nil {
		ex.infoLog(msg)
	}
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
	switch f := expr.Fun.(type) {
	case *ast.SelectorExpr:
		switch ft := f.X.(type) {
		case *ast.CallExpr: // path.call(xxx).LocaleString(p)
			return ex.inspect(ft, info)
		case *ast.Ident: // path.call(xxx) 或是 path.Type.call(xxx) 或是 Type.call(xxx)
			obj := info.ObjectOf(ft)
			if obj == nil {
				break
			}

			switch o := obj.(type) {
			case *types.PkgName:
				if t := info.ObjectOf(f.Sel); t != nil { // 可能指向其它包的别名
					if tn, ok := t.(*types.TypeName); ok && tn.IsAlias() {
						pkgName, structName := getTypeName(tn.Type().String())
						if !ex.tryAppendMsg(expr, pkgName, "", structName) {
							return false
						}
					}
				}

				pkgName := o.Imported().Path()
				return ex.tryAppendMsg(expr, pkgName, "", f.Sel.Name)
			case *types.Var, *types.Const, *types.Nil:
				pkgName, structName := getTypeName(o.Type().String())
				return ex.tryAppendMsg(expr, pkgName, structName, f.Sel.Name)
			default: // 其它可能类型：Func
				pos := ex.fset.Position(ft.Pos())
				panic(fmt.Sprintf("未正确处理 %T 类型的对象,位于 %s", o, pos))
			}
		default:
			return true
		}
	case *ast.Ident: // call(xxx) 调用当前包中的函数或是类型转换，肯定不会有结构体相关联。
		if obj := info.ObjectOf(f); obj != nil {
			if tn, ok := obj.(*types.TypeName); ok && tn.IsAlias() {
				pkgName, structName := getTypeName(tn.Type().String())
				if !ex.tryAppendMsg(expr, pkgName, "", structName) {
					return false
				}
			}

			var pkgName string
			if pkg := obj.Pkg(); pkg != nil {
				pkgName = pkg.Path()
			}
			return ex.tryAppendMsg(expr, pkgName, "", f.Name)
		}
	}
	return true
}

func (ex *extractor) tryAppendMsg(expr *ast.CallExpr, pkgName, structName, name string) (continueInspect bool) {
	exists := sliceutil.Exists(ex.funcs, func(m fn, _ int) bool {
		return m.name == name && pkgName == m.pkgName && structName == m.typeName
	})
	if !exists {
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

	if sliceutil.Exists(ex.msg, func(m message.Message, _ int) bool { return m.Key == key }) {
		ex.warnLog(localeutil.Phrase("has same key %s at %s:%d, will be ignore", strconv.Quote(key), path, p.Line))
		return
	}

	ex.info(localeutil.Phrase("find new locale string %s at %s:%d", strconv.Quote(key), path, p.Line))
	ex.msg = append(ex.msg, message.Message{Key: key, Message: message.Text{Msg: key}})
}

func getTypeName(t string) (pkg, structure string) {
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
