package app

import (
	"context"
	"fmt"
)

type (
	App interface {
		Hello(ctx context.Context, msg string) (string, error)
	}

	app struct{}
)

// Hello implements App.
func (*app) Hello(ctx context.Context, msg string) (string, error) {
	return fmt.Sprintf("Hello %s", msg), nil
}

var _ (App) = (*app)(nil)

func New() (*app, error) {
	return &app{}, nil
}
