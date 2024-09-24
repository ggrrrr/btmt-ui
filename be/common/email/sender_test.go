package email

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	cfg        = Config{}
	email_from = ""
	email1     = ""
	email2     = ""
	repoFolder = ""
)

func loadConfig() {
	email_from = os.Getenv("EMAIL_FROM1")
	email1 = os.Getenv("EMAIL_EMAIL1")
	email2 = os.Getenv("EMAIL_EMAIL2")
	repoFolder = os.Getenv("REPO_FOLDER")
	cfg = Config{
		SMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
		SMTPAddr: os.Getenv("EMAIL_SMTP_ADDR"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}

func TestDialAndSend(t *testing.T) {
	newLine = []byte("\r\n")

	loadConfig()
	t.Skip("NO Addr CONFIG")
	if cfg.SMTPAddr == "" {
		t.Skip("NO Addr CONFIG")
	}
	var err error
	type Data struct {
		User string
	}
	myData := Data{User: "Pesho"}
	template_data := `<p>Скъпи <b>{{ .User }}</b>, welcome to <img src="cid:glass-mug-variant.png" alt="My image" /></p>.`
	tmpl := template.Must(template.New("template_data").Parse(template_data))

	email, err := CreateMsg(
		Rcpt{Mail: email_from, Name: "admin at batamata.org"},
		[]Rcpt{{Mail: email1, Name: "Vesko"}},
		"проба ?{}<> с символи!",
	)
	require.NoError(t, err)

	// email.AddBcc(RcptList{{Mail: "mandajiev@yahoo.com", Name: "Besko"}})
	email.AddCc(RcptList{{Mail: email2, Name: "Besko"}})
	email.AddFile(fmt.Sprintf("%s/glass-mug-variant.png", repoFolder))
	email.AddHtmlBodyWriter(func(w io.Writer) error {
		return tmpl.Execute(w, myData)
	})

	client, err := NewSender(cfg)
	require.NoError(t, err)
	defer client.Close()

	err = client.Send(email)
	assert.NoError(t, err)
}

func TestMultipleMsg(t *testing.T) {
	newLine = []byte("\r\n")
	loadConfig()
	t.Skip("NO Addr CONFIG")
	if cfg.SMTPAddr == "" {
		t.Skip("NO Addr CONFIG")
	}
	emails := RcptList{
		Rcpt{Mail: email1, Name: "email 1"},
		Rcpt{Mail: email2, Name: "email 2"},
	}
	type Data struct {
		Time time.Time
		Name string
	}
	template_data := `<p>Hello: <b>{{ .Name }}</b>, time is: {{ .Time }} /></p>.`
	tmpl := template.Must(template.New("template_data").Parse(template_data))

	conn, err := NewSender(cfg)
	require.NoError(t, err)
	defer conn.Close()

	for _, m := range emails {
		data := Data{
			Name: m.Name,
			Time: time.Now(),
		}
		msg, err := CreateMsg(
			Rcpt{Mail: email_from, Name: "admin at batamata.org"},
			[]Rcpt{m},
			"testing mails",
		)
		require.NoError(t, err)
		msg.AddHtmlBodyWriter(func(w io.Writer) error {
			return tmpl.Execute(w, data)
		})
		err = conn.Send(msg)
		assert.NoError(t, err)

	}

}

func TestAuth(t *testing.T) {
	newLine = []byte("\r\n")

	testSender := &Sender{
		cfg:        Config{},
		tcpConn:    nil,
		smtpClient: nil,
	}

	tests := []struct {
		name      string
		prep      func()
		respErr   error
		respErrAs error
	}{
		{
			name: "ok",
			prep: func() {
				testSender.smtpClient = NewSmtpClientMock()
			},
			respErr:   nil,
			respErrAs: nil,
		},
		{
			name: "extension error",
			prep: func() {
				smtpClient := NewSmtpClientMock()
				smtpClient.falseOnExtension = true
				testSender.smtpClient = smtpClient
			},
			respErr:   nil,
			respErrAs: nil,
		},
		{
			name: "starttls error",
			prep: func() {
				smtpClient := NewSmtpClientMock()
				smtpClient.errorOnStartTLS = true
				testSender.smtpClient = smtpClient
			},
			respErr:   fmt.Errorf("starttls"),
			respErrAs: &SmtpAuthError{},
		},
		{
			name: "auth error",
			prep: func() {
				smtpClient := NewSmtpClientMock()
				smtpClient.errorOnAuth = true
				testSender.smtpClient = smtpClient
			},
			respErr:   fmt.Errorf("auth"),
			respErrAs: &SmtpAuthError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prep()
			err := testSender.smtpAuth()
			if tc.respErr != nil {
				// assert.Equal(t, tc.respErr, err)
				assert.ErrorAs(t, err, &tc.respErr)
				assert.ErrorAs(t, err, &tc.respErrAs)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func TestSend(t *testing.T) {

	newLine = []byte("\n")

	smtpClientMock := NewSmtpClientMock()

	testSender := &Sender{
		cfg: Config{
			SMTPHost: "",
			SMTPAddr: "",
			Username: "",
			Password: "",
			AuthType: "",
			Timeout:  time.Second * 1,
		},
		tcpConn:    nil,
		smtpClient: smtpClientMock,
	}
	email, err := CreateMsg(
		Rcpt{Mail: "mail@from", Name: "name from"},
		[]Rcpt{{Mail: "mail@to", Name: "name to"}},
		"mail subject",
	)
	require.NoError(t, err)
	email.AddBodyString("mail body")

	expectedData := `From: "name from" <mail@from>
To: "name to" <mail@to>
Subject: mail subject
MIME-Version: 1.0
Content-Type: multipart/related; boundary=ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910

--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910
Content-Type: text/plain

mail body
--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910--
`

	err = testSender.Send(email)
	require.NoError(t, err)

	fmt.Printf("\n\n%v\n\n", email.rootWriter.boundary)
	fmt.Printf("\n\n%v\n\n", smtpClientMock.dataBlocks)
	testMockedEmail(t, email, expectedData, smtpClientMock.dataBlocks[0])
}

func testMockedEmail(t *testing.T, email *Msg, expectedData string, actualData string) {
	expected := strings.ReplaceAll(expectedData, "ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910", email.rootWriter.boundary)
	require.True(t, actualData != "", "actual data is empty")
	require.Equal(t, expected, actualData, "data dont match")
}
