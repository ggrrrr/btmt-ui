package ver

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func BuildVersion(buildVerFile string) string {
	verF, err := os.ReadFile(buildVerFile)
	if err != nil {
		return "dev"
	}
	return strings.TrimSpace(string(verF))
}

func BuildTime(buildTimeFile string) time.Time {

	if buildTimeFile == "" {
		return time.Now()
	}

	verTS, err := os.ReadFile(buildTimeFile)
	if err != nil {
		fmt.Printf("cant read file name:%s  %#v\n", buildTimeFile, err)
		return time.Now().UTC()
	}

	buildTs, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(string(verTS)))
	if err != nil {
		fmt.Printf("file name:%s from string:[%s] parsing error: %#v\n", buildTimeFile, string(verTS), err)
		return time.Now().UTC()
	}
	return buildTs
}
