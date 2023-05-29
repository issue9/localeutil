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
func ReadDomains(key string, domain ...string) (string, error) {
	for _, d := range domain {
		if l, err := Read(key, d); err == nil && l != "" {
			return l, nil
		}
	}
	return "", nil
}

func Read(key, domain string) (string, error) {
	b := &bytes.Buffer{}

	cmd := exec.Command("defaults", "read", domain, key)
	cmd.Stdout = b
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(b.String()), nil
}

func Write(domain, key, t string, value ...string) error {
	if t[0] != '-' {
		t = "-" + t
	}

	args := []string{"write", domain, key, t}
	args = append(args, value...)
	return exec.Command("defaults", args...).Run()
}
