package main

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

func ListPaths(root string) ([]string, error) {
	var paths []string
	if err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %q: %s", p, err)
		}
		p = path.Join("/", filepath.ToSlash(rel))
		if d.IsDir() && !strings.HasSuffix(p, "/") {
			p += "/"
		}
		paths = append(paths, p)
		return nil
	}); err != nil {
		return nil, err
	}
	return paths, nil
}
