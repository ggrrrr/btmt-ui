package system

import (
	"github.com/ggrrrr/btmt-ui/be/common/postgres"
)

func (s *System) initDB() error {
	if s.cfg.Postgres.Host == "" {
		return nil
	}
	db, err := postgres.Connect(s.cfg.Postgres)
	if err != nil {
		return err
	}
	s.db = db
	// s.waiter.Cleanup()
	s.waiter.Cleanup(func() {
		db.Close()
	})
	return nil
}
