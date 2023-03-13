// SPDX-License-Identifier: MIT

package syslocale

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func getLocaleName() (string, error) {
	if l := getEnvLang(); l != "" { // 优先看环境变量
		return l, nil
	}

	if l, err := parseDefaultAppleLocale("-g"); err == nil {
		return l, nil
	}
	return parseDefaultAppleLocale("/Library/Preferences/.GlobalPreferences")
}

func parseDefaultAppleLocale(t string) (string, error) {
	b := &bytes.Buffer{}

	cmd := exec.Command("defaults", "read", t, "AppleLocale")
	cmd.Stdout = b
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("检测用户环境时返回错误：%w", err)
	}

	return strings.TrimSpace(b.String()), nil
}
