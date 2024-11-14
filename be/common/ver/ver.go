package ver

import (
	_ "embed"
	"strings"
	"time"
)

//go:embed build_ver.txt
var buildVer string

//go:embed build_ts.txt
var buildTs string

func BuildVersion() string {
	return buildVer
}

func BuildTime() time.Time {
	ts2, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(buildTs))
	if err != nil {
		// fmt.Println("BuildTime error: ", err)
		return time.Now().UTC()
	}
	return ts2
}
