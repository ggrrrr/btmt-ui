package ver

import (
	_ "embed"
	"strings"
	"time"
)

//go:embed ver.txt
var v string

//go:embed ts.txt
var ts string

func BuildVersion() string {
	return v
}

func BuildTime() time.Time {
	ts2, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(ts))
	if err != nil {
		// fmt.Println("BuildTime error: ", err)
		return time.Now().UTC()
	}
	return ts2
}
