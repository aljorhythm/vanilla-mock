package test

import (
	"path"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _   = runtime.Caller(0)
	d            = path.Join(path.Dir(b))
	project_root = filepath.Dir(d)
)

func GetProjectRoot() string {
	return project_root
}

func GetTestRoot() string {
	return d
}

func GetAbsPath(relativePath string) string {
	return path.Join(GetProjectRoot(), relativePath)
}
