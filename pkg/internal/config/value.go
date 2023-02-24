//go:build !windows

package config

import (
	"os"
	"path/filepath"
)

func GetRoot() string {
	root := filepath.Join("/opt", "gpm")
	_ = os.MkdirAll(root, os.ModePerm)
	return root
}

var DefaultRoot = GetRoot()
