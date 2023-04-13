// SPDX-License-Identifier: MIT

// Package syslocale 获取所在系统的本地化语言信息
package syslocale

import (
	"os"
	"strings"
)

// Get 返回当前系统的本地化信息
func Get() string {
	if lang := getEnv("LANGUAGE"); lang != "" {
		return lang
	}

	if lang := getOSLocaleName(); lang != "" {
		return trimLocale(lang)
	}

	for _, env := range [...]string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if lang := getEnv(env); lang != "" {
			return lang
		}
	}

	return ""
}

func getEnv(env string) string { return trimLocale(os.Getenv(env)) }

// zh_CN.UTF-8 过滤掉最后的编码方式
func trimLocale(v string) string {
	if index := strings.LastIndexByte(v, '.'); index > 0 {
		v = v[:index]
	}
	return strings.TrimSpace(v)
}
