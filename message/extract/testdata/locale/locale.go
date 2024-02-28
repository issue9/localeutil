// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package locale

// 测试包内的各类型数据

type (
	String string

	Printer struct{}

	Alias = Printer

	Interface interface {
		Printf(string) string
	}

	interfaceImpl struct{}

	GPrinter[T any] struct{}

	IntAlias = GPrinter[int] // 在 type() 中的类型声明，调用时没有关联的 Object
)

type StrAlias = GPrinter[string]

var _ Interface = &interfaceImpl{}

func Print(string) string { return "" }

func (*Printer) Print(string) string { return "" }

func (*GPrinter[T]) GPrint(string) string { return "" }

func (*interfaceImpl) Printf(string) string { return "" }

func NewPrinter() *Printer { return &Printer{} }
