// SPDX-License-Identifier: MIT

package testdata

// 测试当前包的引用

func Print(key string) {}

type Printer struct{}

func (p *Printer) Print(key string) {}

func local() {
	Print("local print")

	var p *Printer = &Printer{}
	p.Print("local struct print")
}
