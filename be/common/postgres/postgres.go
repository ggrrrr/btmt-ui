package postgres

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type (
	Config struct {
		Host     string `env:"HOST"`
		Port     int    `env:"PORT" envDefault:"5432"`
		Username string `env:"USERNAME"`
		Password string `env:"PASSWORD"`
		Database string `env:"DATABASE"`
		SSLMode  string `env:"SSL_MODE"`
		Prefix   string `env:"PREFIX"`
		Schema   string `env:"SCHEMA"`
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
		logger.Error(err).
			Str("host", cfg.Host).
			Str("database", cfg.Database).
			Msg("Ping")
		return nil, err
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "dev"
	}
	logger.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("User", cfg.Username).
		Str("Database", cfg.Database).
		Msg("Connected")

	return db, nil
}
