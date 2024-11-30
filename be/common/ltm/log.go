package ltm

import (
	"log/slog"
	"os"
)

func New() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	slog.SetDefault(logger)
}
