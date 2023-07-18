// SPDX-License-Identifier: MIT

package extract

import (
	"io/fs"
	"path/filepath"
)

func getDir(root string, r bool) ([]string, error) {
	if !r {
		return []string{root}, nil
	}

	dirs := make([]string, 0, 30)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err == nil && d.IsDir() {
			dirs = append(dirs, p)
		}
		return err
	})

	if err != nil {
		return nil, err
	}
	return dirs, nil
}
