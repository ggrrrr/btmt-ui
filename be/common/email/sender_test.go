package email

import (
	"html/template"
	"io"
	"os"
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
)

func loadConfig() {
	email_from = os.Getenv("EMAIL_FROM1")
	email1 = os.Getenv("EMAIL_EMAIL1")
	email2 = os.Getenv("EMAIL_EMAIL2")
	cfg = Config{
		Host:     os.Getenv("EMAIL_HOST"),
		Addr:     os.Getenv("EMAIL_ADDR"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}

func TestDialAndSend(t *testing.T) {
	loadConfig()
	if cfg.Addr == "" {
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
	email.AddFile("/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/glass-mug-variant.png")
	email.AddHtmlBodyWriter(func(w io.Writer) error {
		return tmpl.Execute(w, myData)
	})

	client, err := Dial(cfg)
	require.NoError(t, err)
	defer client.Close()

	err = client.Send(email)
	assert.NoError(t, err)
}

func TestMultipleMsg(t *testing.T) {
	loadConfig()

	if cfg.Addr == "" {
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

	conn, err := Dial(cfg)
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
