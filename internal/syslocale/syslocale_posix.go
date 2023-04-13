// SPDX-License-Identifier: MIT

//go:build !windows && !js && !darwin
// +build !windows,!js,!darwin

package syslocale

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 依次读取以下几个路径
//   - XDG_CONFIG/locale.conf
//   - HOME/locale.conf
//   - /etc/locale.conf
func getOSLocaleName() string {
	if dir, err := os.UserConfigDir(); err == nil {
		if val := readFromFile(dir); val != "" {
			return val
		}
	}

	if dir, err := os.UserHomeDir(); err == nil {
		if val := readFromFile(dir); val != "" {
			return val
		}
	}

	if val := readFromFile("/etc"); val != "" {
		return val
	}

	return ""
}

func readFromFile(dir string) string {
	f, err := os.Open(filepath.Join(dir, "locale.conf"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Println(err)
		}
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		vals := strings.Split(scanner.Text(), "=")
		if len(vals) != 2 {
			continue
		}
		if vals[0] != "LANG" && vals[0] != "LC_ALL" && vals[0] != "LC_MESSAGES" {
			continue
		}

		if val := strings.TrimSpace(vals[1]); val != "" {
			return strings.Trim(val, `"`)
		}
	}
	return ""
}
