package email

import "context"

type (
	MockSmtpConnector struct {
		ForErr error
		Sender *MockSmtpSender
	}

	MockSmtpSender struct {
		ForErr   error
		LastMail string
	}
)

var _ (EmailConnector) = (*MockSmtpConnector)(nil)

var _ (EmailSender) = (*MockSmtpSender)(nil)

// Connect implements Connector.
func (m *MockSmtpConnector) Connect(ctx context.Context) (EmailSender, error) {
	if m.ForErr != nil {
		return nil, m.ForErr
	}
	return m.Sender, nil
}

// Close implements Sender.
func (m *MockSmtpSender) Close() error {
	return nil
}

// Send implements Sender.
func (m *MockSmtpSender) Send(ctx context.Context, email *Msg) error {
	if m.ForErr != nil {
		return m.ForErr
	}

	m.LastMail = email.DumpToText()
	return nil
}
