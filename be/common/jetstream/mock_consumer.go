package jetstream

import "context"

type MockConsumer struct {
	Err         error
	HandlerFunc DataHandlerFunc
}

type MockConnecter struct {
	Err         error
	HandlerFunc DataHandlerFunc
}

var _ (Consumer) = (*MockConsumer)(nil)

var _ (Consumer) = (*MockConsumer)(nil)

// Consume implements Consumer.
func (m *MockConsumer) Consume(ctx context.Context, handler DataHandlerFunc) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}

// Shutdown implements Consumer.
func (m *MockConsumer) Shutdown() error {
	return nil
}
