package help

import (
	"fmt"
	"runtime"
	"strings"
)

func RepoDir() string {
	return dir()
}
func dir() string {
	_, filename, _, _ := runtime.Caller(1)
	files := strings.Split(filename, "/")
	fmt.Printf("help/repodir.go: %v\n", filename)
	rootDir := strings.Join(files[:len(files)-3], "/")
	return rootDir
}
