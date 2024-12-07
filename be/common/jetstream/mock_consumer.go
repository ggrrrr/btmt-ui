package jetstream

import "context"

type MockConsumer struct {
	Err         error
	HandlerFunc DataHandler
}

type MockConnecter struct {
	Err         error
	HandlerFunc DataHandler
}

var _ (Consumer) = (*MockConsumer)(nil)

var _ (Consumer) = (*MockConsumer)(nil)

// Consume implements Consumer.
func (m *MockConsumer) Consume(ctx context.Context, handler DataHandler) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}

// Shutdown implements Consumer.
func (m *MockConsumer) Shutdown() error {
	return nil
}
