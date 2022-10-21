package test

import (
	"path"
	"runtime"
)

func MountAbsolutPath(p string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), p)
}
