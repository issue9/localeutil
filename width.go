// SPDX-License-Identifier: MIT

package localeutil

import "golang.org/x/text/width"

var defaultWidthOptions = WidthOptions{
	width.EastAsianFullwidth: 2,
	width.EastAsianWide:      2,

	width.EastAsianHalfwidth: 1,
	width.EastAsianNarrow:    1,

	width.Neutral:            1,
	width.EastAsianAmbiguous: 1,
}

// WidthOptions 用于指定各类字符的宽度
//
// 拥有以下几个配置项：
//
//	width.EastAsianFullwidth: 2
//	width.EastAsianWide:      2
//	width.EastAsianHalfwidth: 1
//	width.EastAsianNarrow:    1
//	width.Neutral:            1
//	width.EastAsianAmbiguous: 1
//
// 对于 [width.EastAsianAmbiguous] 不同的字体可能有不同的设置。
type WidthOptions map[width.Kind]int

// DefaultWidthOptions 返回默认的配置项的副本
//
// 如果要基于默认值作修改，可以采用此方法。
func DefaultWidthOptions() WidthOptions {
	o := make(WidthOptions, len(defaultWidthOptions))
	for k, v := range defaultWidthOptions {
		o[k] = v
	}
	return o
}

// Width 计算字符串的宽度
func (wo WidthOptions) Width(s string) (w int) {
	for _, r := range s {
		w += wo[width.LookupRune(r).Kind()]
	}
	return w
}

// Width 采用 defaultWidthOptions 计算字符串的宽度
//
// 如果有特殊要求，可以使用 [WidthOptions] 自定义各类字符的宽度。
func Width(s string) int { return defaultWidthOptions.Width(s) }
