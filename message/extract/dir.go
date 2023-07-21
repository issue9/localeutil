// SPDX-License-Identifier: MIT

package extract

import (
	"io/fs"
	"os"
	"path/filepath"
)

func getDir(root string, r, skip bool) ([]string, error) {
	if !r {
		return []string{root}, nil
	}

	dirs := make([]string, 0, 30)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		if skip && root != p {
			stat, err := os.Stat(filepath.Join(p, "go.mod"))
			if err == nil && !stat.IsDir() {
				return fs.SkipDir
			}
		}

		dirs = append(dirs, p)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return dirs, nil
}
