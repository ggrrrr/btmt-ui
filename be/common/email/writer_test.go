package email

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteTo(t *testing.T) {
	pwd := os.Getenv("PWD")
	newLine = []byte("\n")
	mail, _ := CreateMsg(
		Rcpt{Mail: "asd@asd", Name: "From name"},
		RcptList{Rcpt{Mail: "to@to", Name: "To name"}},
		"Subject 1",
	)

	result := `Subject: Subject 1
From: "From name" <asd@asd>
To: "To name" <to@to>
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
	mail.AddFile(fmt.Sprintf("%s/../../../test.txt", pwd))

	buf := new(bytes.Buffer)
	err := mail.writerTo(buf)
	assert.NoError(t, err)
	assert.True(t, mail.rootWriter.boundary != "")
	resultB := strings.ReplaceAll(result, "c306cfe3a93593b66e74fae51429f45917f23e2a74f3e70aa042fdec2891", mail.rootWriter.boundary)

	fmt.Printf("%s \n", buf.String())
	assert.Equal(t, resultB, buf.String())
}
