// SPDX-License-Identifier: MIT

package localeutil

import "golang.org/x/text/width"

// Width 计算字符串的宽度
func Width(s string) (w int) {
	for _, r := range s {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			w += 2
		case width.EastAsianHalfwidth, width.EastAsianNarrow,
			width.Neutral, width.EastAsianAmbiguous:
			w += 1
		}
	}

	return w
}
