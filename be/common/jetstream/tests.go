package jetstream

import "github.com/ggrrrr/btmt-ui/be/common/token"

func ConnectForTest() (*NatsConnection, error) {
	verifier := token.NewVerifierMock()
	var cfg = Config{
		URL: "localhost:4222",
	}

	return Connect(cfg, WithVerifier(verifier))
}
