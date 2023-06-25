// SPDX-License-Identifier: MIT

//go:build darwin || ios

// Package defaults 提供苹果系统功能
package defaults

import (
	"bytes"
	"os/exec"
	"strings"
)

// ReadDomains 从 domain 中查找 key 值
//
// 按顺序找到第一个为止。
func ReadDomains(key string, domain ...string) string {
	for _, d := range domain {
		if l := Read(key, d); l != "" {
			return l
		}
	}
	return ""
}

func Read(key, domain string) string {
	b := &bytes.Buffer{}

	cmd := exec.Command("defaults", "read", domain, key)
	cmd.Stdout = b
	if cmd.Run() != nil {
		// Run() 返回的 err 本身无意思，错误信息存在 cmd.Stderr 或是 cmd.Stdout。
		// 也无法判断是找不到 key 还是 domain，干脆不作其它处理。
		return ""
	}

	return strings.TrimSpace(b.String())
}

func Write(domain, key, t string, value ...string) error {
	if t[0] != '-' {
		t = "-" + t
	}

	args := []string{"write", domain, key, t}
	args = append(args, value...)
	return exec.Command("defaults", args...).Run()
}
