// SPDX-License-Identifier: MIT

package locale

import (
	"fmt"
	"log"
)

// 以下为测试数据

var (
	p1          = NewPrinter() // 需要类型推导
	p2 *Printer = &Printer{}

	i1           = &interfaceImpl{} // 需要判断是否实现了接口，无法实现。
	i2 Interface = &interfaceImpl{}

	a1        = &Alias{} // 需要类型推导
	a2 *Alias = &Alias{}

	g1                 = &GPrinter[uint]{}
	g2 *GPrinter[uint] = &GPrinter[uint]{}

	sg1           = &StrAlias{}
	sg2 *StrAlias = &StrAlias{}

	ig1           = &IntAlias{}
	ig2 *IntAlias = &IntAlias{}
)

const (
	printKey  = "Print key"
	stringKey = "String key"
)

const c = String(stringKey)

var (
	v = Print(printKey)

	_ = p1.Print("p1.Print")
	_ = p2.Print("p2.Print")
	_ = (&Printer{}).Print("p3.Print")

	_ = i1.Printf("i1.Printf") // 不会主动推断接口实现
	_ = i2.Printf("i2.Printf")

	_ = a1.Print("a1.Print")
	_ = a2.Print("a2.Print")

	_ = g1.GPrint("g1.GPrint")
	_ = g2.GPrint("g2.GPrint")

	_ = sg1.GPrint("sg1.GPrint")
	_ = sg2.GPrint("sg2.GPrint")

	_ = ig1.GPrint("ig1.GPrint")
	_ = ig2.GPrint("ig2.GPrint")
)

func output(p *Printer, i1 *interfaceImpl, i2 Interface, a *Alias, g *GPrinter[uint], sg *StrAlias, ig *IntAlias) {
	const c = String(stringKey)
	log.Print(c)

	_ = Print(printKey)

	_ = p.Print("p2.Print") // 由函数参数获得类型信息

	_ = i1.Printf("i1.Printf") // 不会主动判断是否实现了接口
	_ = i2.Printf("i2.Printf")

	_ = a.Print("a2.Print")

	_ = g.GPrint("g2.GPrint")

	_ = sg.GPrint("sg2.GPrint")

	_ = ig.GPrint("ig2.GPrint")

	fmt.Println("abc")
}
