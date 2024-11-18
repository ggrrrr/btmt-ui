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
		return "dev"
	}

	logger.Info().Str("build.version", string(verF)).Send()
	return strings.TrimSpace(string(verF))
}

func BuildTime(buildTimeFile string) time.Time {

	if buildTimeFile == "" {
		return time.Now().UTC()
	}

	verTS, err := os.ReadFile(buildTimeFile)
	if err != nil {
		logger.Error(err).Str("file", buildTimeFile).Msg("build version: cant read file")
		return time.Now().UTC()
	}

	buildTs, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(string(verTS)))
	if err != nil {
		logger.Error(err).Str("file", buildTimeFile).Msg("build version: cant parse time")
		return time.Now().UTC()
	}
	logger.Info().Time("build.time", buildTs).Send()
	return buildTs
}
