// SPDX-License-Identifier: MIT

// Package syslocale 获取所在系统的本地化语言信息
package syslocale

import (
	"os"
	"strings"

	"golang.org/x/text/language"
)

// Get 返回当前系统的本地化信息
func Get() (language.Tag, error) {
	name, err := getLocaleName()
	if err != nil {
		return language.Und, err
	}
	return language.Parse(name)
}

// 获取环境变量 LANG 中有关本地化信息的值。
// https://www.gnu.org/software/gettext/manual/html_node/Locale-Environment-Variables.html
func getEnvLang() string {
	for _, env := range [...]string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if name := os.Getenv(env); len(name) > 0 {
			// zh_CN.UTF-8 过滤掉最后的编码方式
			if index := strings.LastIndexByte(name, '.'); index > 0 {
				name = name[:index]
			}
			return name
		}
	}
	return ""
}
