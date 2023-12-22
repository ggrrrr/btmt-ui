package email

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteTo(t *testing.T) {
	t.Skip("TODO")
	newLine = []byte("\n")
	mail, _ := CreateMsg(
		Rcpt{Mail: "asd@asd", Name: "Name"},
		RcptList{Rcpt{Mail: "to@to", Name: "to"}},
		"Subject 1",
	)

	result := `Subject: Subject 1
From: "Name" <asd@asd>
To: to@to
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
--c306cfe3a93593b66e74fae51429f45917f23e2a74f3e70aa042fdec2891--`
	mail.AddHtmlBodyString("asd")
	mail.AddFile("/Users/vesko/go/src/github.com/ggrrrr/btmt-ui/test.txt")

	buf := new(bytes.Buffer)
	mail.writerTo(buf)
	fmt.Printf("%s \n", buf.String())
	assert.Equal(t, result, buf.String())
}
