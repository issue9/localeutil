// SPDX-License-Identifier: MIT

package syslocale

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func getOSLocaleName() string {
	if l := parseDefaultAppleLocale("-g"); l != "" {
		return l
	}
	return parseDefaultAppleLocale("/Library/Preferences/.GlobalPreferences")
}

func parseDefaultAppleLocale(t string) string {
	b := &bytes.Buffer{}

	cmd := exec.Command("defaults", "read", t, "AppleLocale")
	cmd.Stdout = b
	if err := cmd.Run(); err != nil {
		log.Println(fmt.Errorf("检测用户环境时返回错误：%w", err))
		return ""
	}

	return strings.TrimSpace(b.String())
}
