package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
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
		log.Log().Error(err, "Connect")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Log().Error(err, "Ping",
			slog.String("host", cfg.Host),
			slog.String("database", cfg.Database))
		return nil, err
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "dev"
	}
	log.Log().Error(err, "Connected",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("User", cfg.Username),
		slog.String("Database", cfg.Database))

	return db, nil
}
