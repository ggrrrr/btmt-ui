package email

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/help"
	"github.com/stretchr/testify/assert"
)

var pwd = help.RepoDir()

func Test_WriteTo(t *testing.T) {
	newLine = []byte("\n")
	var err error

	mailEmpty, _ := createMsg(
		Rcpt{addr: "asd@asd", name: "From name"},
		RcptList{Rcpt{addr: "to@to", name: "To name"}},
		"Subject 1",
	)

	bufEmpty := new(bytes.Buffer)
	err = mailEmpty.writerTo(bufEmpty)
	err1 := &MailFormatError{}
	assert.ErrorAs(t, err, &err1)

	mail, _ := createMsg(
		Rcpt{addr: "asd@asd", name: "From name"},
		RcptList{Rcpt{addr: "to@to", name: "To name"}},
		"Subject 1",
	)

	result := `From: "From name" <asd@asd>
To: "To name" <to@to>
Subject: Subject 1
MIME-Version: 1.0
Content-Type: multipart/related; boundary=c306cfe3a93593b66e74fae51429f45917f23e2a74f3e70aa042fdec2891

--c306cfe3a93593b66e74fae51429f45917f23e2a74f3e70aa042fdec2891
Content-Type: text/html

asd
--c306cfe3a93593b66e74fae51429f45917f23e2a74f3e70aa042fdec2891
Content-Type: text/plain; charset=utf-8; name="test.txt"
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="test.txt"
Content-ID: <test.txt>

MQoyCg==
--c306cfe3a93593b66e74fae51429f45917f23e2a74f3e70aa042fdec2891--
`
	mail.AddHtmlBodyString("asd")
	mail.AddFile(fmt.Sprintf("%s/test.txt", pwd))

	buf := new(bytes.Buffer)
	err = mail.writerTo(buf)
	assert.NoError(t, err)
	assert.True(t, mail.rootWriter.boundary != "")
	resultB := strings.ReplaceAll(result, "c306cfe3a93593b66e74fae51429f45917f23e2a74f3e70aa042fdec2891", mail.rootWriter.boundary)

	fmt.Printf("%s \n", buf.String())
	assert.Equal(t, resultB, buf.String())
}
