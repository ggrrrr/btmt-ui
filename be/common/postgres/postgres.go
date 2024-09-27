package postgres

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type (
	Config struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
		SSLMode  string
		Prefix   string
	}
)

func Connect(cfg Config) (*sql.DB, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)

	db, err := otelsql.Open("postgres", psqlConn, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))

	// db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		logger.Error(err).Msg("Connect")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.Error(err).Msg("Ping")
		return nil, err
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "dev"
	}
	logger.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("User", cfg.Username).
		Msg("Connected")

	return db, nil
}
