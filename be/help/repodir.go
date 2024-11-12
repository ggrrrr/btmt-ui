package help

import (
	"runtime"
	"strings"
)

// Returns full path of the repo folder on your local machine
// so unit tests can use is as file root dir
func RepoDir() string {
	return dir()
}

func dir() string {
	_, filename, _, _ := runtime.Caller(1)
	files := strings.Split(filename, "/")
	rootDir := strings.Join(files[:len(files)-3], "/")
	return rootDir
}
