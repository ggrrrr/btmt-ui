package help

import (
	"runtime"
	"strings"
)

func RepoDir() string {
	return dir()
}
func dir() string {
	_, filename, _, _ := runtime.Caller(1)
	files := strings.Split(filename, "/")
	rootDir := strings.Join(files[:len(files)-3], "/")
	return rootDir
}
