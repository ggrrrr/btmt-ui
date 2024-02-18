package postgres

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

