package email

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/help"
)

var (
	cfg        = Config{}
	email_from = ""
	email1     = ""
	email2     = ""
	repoFolder = help.RepoDir()
)

func LoadTestCfg() {
	email_from = os.Getenv("EMAIL_FROM1")
	email1 = os.Getenv("EMAIL_EMAIL1")
	email2 = os.Getenv("EMAIL_EMAIL2")
	cfg = Config{
		SMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
		SMTPAddr: os.Getenv("EMAIL_SMTP_ADDR"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}

}

func TestIntDialAndSend(t *testing.T) {
	newLine = []byte("\r\n")

	LoadTestCfg()
	t.Skip("NO Addr CONFIG")
	ctx := context.Background()
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

	email, err := createMsg(
		Rcpt{addr: email_from, name: "admin at batamata.org"},
		[]Rcpt{{addr: email1, name: "Vesko"}},
		// "проба ?{}<> с символи!",
	)
	require.NoError(t, err)

	// email.AddBcc(RcptList{{Mail: "mandajiev@yahoo.com", Name: "Besko"}})
	email.AddCc(RcptList{{addr: email2, name: "Besko"}})
	email.AddFile(fmt.Sprintf("%s/glass-mug-variant.png", repoFolder))
	email.AddHtmlBodyWriter(func(w io.Writer) error {
		return tmpl.Execute(w, myData)
	})

	client, err := cfg.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	err = client.Send(ctx, email)
	assert.NoError(t, err)

}

func TestIntMultipleMsg(t *testing.T) {

	newLine = []byte("\r\n")
	LoadTestCfg()
	t.Skip("NO Addr CONFIG")
	if cfg.SMTPAddr == "" {
		t.Skip("NO Addr CONFIG")
	}
	emails := RcptList{
		Rcpt{addr: email1, name: "email 1"},
		Rcpt{addr: email2, name: "email 2"},
	}
	type Data struct {
		Time time.Time
		Name string
	}
	template_data := `<p>Hello: <b>{{ .Name }}</b>, time is: {{ .Time }} /></p>.`
	tmpl := template.Must(template.New("template_data").Parse(template_data))

	ctx := context.Background()
	conn, err := cfg.Connect(ctx)
	require.NoError(t, err)
	defer conn.Close()

	for _, m := range emails {
		data := Data{
			Name: m.name,
			Time: time.Now(),
		}
		msg, err := createMsg(
			Rcpt{addr: email_from, name: "admin at batamata.org"},
			[]Rcpt{m},
			// "testing mails",
		)
		require.NoError(t, err)
		msg.AddHtmlBodyWriter(func(w io.Writer) error {
			return tmpl.Execute(w, data)
		})
		err = conn.Send(ctx, msg)
		assert.NoError(t, err)

	}

}

func TestAuth(t *testing.T) {
	newLine = []byte("\r\n")

	testSender := &sender{
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
				testSender.smtpClient = newSmtpClientMock()
			},
			respErr:   nil,
			respErrAs: nil,
		},
		{
			name: "extension error",
			prep: func() {
				smtpClient := newSmtpClientMock()
				smtpClient.falseOnExtension = true
				testSender.smtpClient = smtpClient
			},
			respErr:   nil,
			respErrAs: nil,
		},
		{
			name: "starttls error",
			prep: func() {
				smtpClient := newSmtpClientMock()
				smtpClient.errorOnStartTLS = true
				testSender.smtpClient = smtpClient
			},
			respErr:   fmt.Errorf("starttls"),
			respErrAs: &SmtpAuthError{},
		},
		{
			name: "auth error",
			prep: func() {
				smtpClient := newSmtpClientMock()
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
				// fmt.Printf("error: %v\n", err)
				assert.ErrorAs(t, err, &tc.respErr)
				assert.ErrorAs(t, err, &tc.respErrAs)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func TestSend(t *testing.T) {
	ctx := context.Background()
	newLine = []byte("\n")

	var testMsg *Msg

	tests := []struct {
		name      string
		prep      func(t *testing.T)
		dataBlock string
		sendErrAs error
	}{
		{
			name: "error empty body",
			prep: func(t *testing.T) {
				var err error
				testMsg, err = createMsg(
					Rcpt{addr: "mail@from", name: "name from"},
					[]Rcpt{{addr: "mail@to", name: "name to"}},
					// "mail subject",
				)
				require.NoError(t, err)
				testMsg.AddCc(RcptList{Rcpt{addr: "cc@cc.cc", name: "cc name"}})
				testMsg.AddBcc(RcptList{Rcpt{addr: "bcc@bcc.bcc", name: "bcc name"}})
			},
			dataBlock: ``,
			sendErrAs: &MailFormatError{},
		},
		{
			name: "single text email",
			prep: func(t *testing.T) {
				var err error
				testMsg, err = createMsg(
					Rcpt{addr: "mail@from", name: "name from"},
					[]Rcpt{{addr: "mail@to", name: "name to"}},
					// "mail subject",
				)
				require.NoError(t, err)
				err = testMsg.SetSubject("mail subject")
				require.NoError(t, err)

				testMsg.AddCc(RcptList{Rcpt{addr: "cc@cc.cc", name: "cc name"}})
				testMsg.AddBcc(RcptList{Rcpt{addr: "bcc@bcc.bcc", name: "bcc name"}})
				assert.NoError(t, err, "prep email")
				testMsg.AddBodyString("mail body")
			},
			dataBlock: `From: "name from" <mail@from>
To: "name to" <mail@to>
Subject: mail subject
Cc: "cc name" <cc@cc.cc>
MIME-Version: 1.0
Content-Type: multipart/related; boundary=ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910

--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910
Content-Type: text/plain

mail body
--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910--
`,
		},
		{
			name: "single html email and file",
			prep: func(t *testing.T) {
				var err error
				testMsg, err = createMsg(
					Rcpt{addr: "mail@from", name: "name from"},
					[]Rcpt{{addr: "mail@to", name: "name to"}},
					// "mail subject",
				)
				require.NoError(t, err)
				err = testMsg.SetSubject("mail subject")
				require.NoError(t, err)

				myData := struct{ User string }{User: "Pesho"}
				template_data := `User: {{ .User }}`
				tmpl := template.Must(template.New("template_data").Parse(template_data))

				testMsg.AddFile(fmt.Sprintf("%s/test.txt", repoFolder))
				testMsg.AddAttachment("myfile.png", func(w io.Writer) error {
					fileBody := "secret"
					_, err := w.Write([]byte(fileBody))
					return err
				},
				)
				testMsg.AddHtmlBodyWriter(
					func(w io.Writer) error {
						return tmpl.Execute(w, myData)
					},
				)
			},
			dataBlock: `From: "name from" <mail@from>
To: "name to" <mail@to>
Subject: mail subject
MIME-Version: 1.0
Content-Type: multipart/related; boundary=ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910

--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910
Content-Type: text/html

User: Pesho
--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910
Content-Type: text/plain; charset=utf-8; name="test.txt"
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="test.txt"
Content-ID: <test.txt>

MQoyCg==
--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910
Content-Type: image/png; name="myfile.png"
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="myfile.png"
Content-ID: <myfile.png>

c2VjcmV0
--ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910--
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smtpClientMock := newSmtpClientMock()

			testSender := &sender{
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

			tc.prep(t)
			err := testSender.Send(ctx, testMsg)
			if tc.sendErrAs == nil {
				require.NoError(t, err)
				testMockedEmail(t, testMsg, tc.dataBlock, smtpClientMock)
			} else {
				assert.ErrorAs(t, err, &tc.sendErrAs)
			}
		})
	}

}

func testMockedEmail(t *testing.T, email *Msg, expectedData string, actualData *smtpClientMock) {
	expected := strings.ReplaceAll(expectedData, "ce82d13b7cf05644c1a5c74b4c700dae854b1213f93ddf4fb12d7fb0c910", email.rootWriter.boundary)
	require.True(t, actualData.dataBlocks[0] != "", "actual data is empty")
	assert.Equal(t, expected, actualData.dataBlocks[0], "data dont match")
	assert.Equal(t, email.from.addr, actualData.from, "from dont match")
	assert.Equal(t, strings.Split(email.to.AddressList(), ","), actualData.to, "to dont match")
}

func TestCreateSender(t *testing.T) {
	ctx := context.Background()

	tcpServerMock, err := newTCPServerMock(":12346")
	require.NoError(t, err)
	go tcpServerMock.Start()

	cfg := Config{
		SMTPHost: "localhost",
		SMTPAddr: ":12346",
		Username: "user",
		Password: "pass",
		AuthType: "",
		Timeout:  time.Second * 2,
	}

	smtp, err := cfg.Connect(ctx)
	require.NoError(t, err)

	err = smtp.smtpClient.Quit()
	require.NoError(t, err)
	smtp.Close()

}

func TestCreateSenderError(t *testing.T) {
	ctx := context.Background()

	tcpServerMock, err := newTCPServerMock(":12345")
	require.NoError(t, err)
	go tcpServerMock.Start()

	cfg := Config{
		SMTPHost: "localhost",
		SMTPAddr: ":12345",
		Username: "user",
		Password: "apassasdasd",
		AuthType: "",
		Timeout:  time.Second * 2,
	}

	_, err = cfg.Connect(ctx)
	require.Error(t, err)
	authErr := &SmtpAuthError{}
	assert.ErrorAs(t, err, &authErr)

}
