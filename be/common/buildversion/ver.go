/*

This folder is to enable setting app version during build
Check Dockerfile and be/common/system/system.go

**/

package buildversion

import (
	"os"
	"strings"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

func BuildVersion(buildVerFile string) string {
	verF, err := os.ReadFile(buildVerFile)
	if err != nil {
		logger.Error(err).Msg("BuildVersion")
		return "dev"
	}

	logger.Info().Str("BuildVersion", string(verF)).Send()
	return strings.TrimSpace(string(verF))
}

func BuildTime(buildTimeFile string) time.Time {

	if buildTimeFile == "" {
		return time.Now().UTC()
	}

	verTS, err := os.ReadFile(buildTimeFile)
	if err != nil {
		logger.Error(err).Msg("BuildTime")
		return time.Now().UTC()
	}

	buildTs, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(string(verTS)))
	if err != nil {
		logger.Error(err).Str("file", buildTimeFile).Msg("BuildTime: cant parse time")
		return time.Now().UTC()
	}
	logger.Info().Time("build.time", buildTs).Msg("BuildTime")
	return buildTs
}
