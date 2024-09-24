package email

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrites(t *testing.T) {
	newLine = []byte("\r\n")

	type tst struct {
		name     string
		want     string
		testFunc func(w *partWriter) error
	}
	tests := []tst{
		{
			name: "writeSting",
			want: "string1",
			testFunc: func(w *partWriter) error {
				return w.writeString("string1")
			},
		},
		{
			name: "writeBoundaryStart",
			want: "\r\n--boundary-1\r\n",
			testFunc: func(w *partWriter) error {
				return w.writeBoundaryStart()
			},
		},
		{
			name: "writeBoundaryClose",
			want: "\r\n--boundary-1--\r\n",
			testFunc: func(w *partWriter) error {
				return w.writeBoundaryClose()
			},
		},
		{
			name: "writeHeader",
			want: "Cc: asd0, asd1\r\n",
			testFunc: func(w *partWriter) error {
				return w.writeHeader(smtpHeader{headerCc, []string{"asd0", "asd1"}})
			},
		},
		{
			name: "writePart",
			want: "\r\n--boundary-1\r\nContent-Type: text/plain\r\n\r\nbody of email",
			testFunc: func(w *partWriter) error {
				return w.writePart(&mailPart{
					contentType: contentTypePlain,
					copier: func(w io.Writer) error {
						// nolint: errcheck
						w.Write([]byte(`body of email`))
						return nil
					},
					encoding: Unencoded,
				})
			},
		},
		{
			name: "writeAttachment",
			want: "\r\n--boundary-1\r\nContent-Type: text/plain; charset=utf-8; name=\"file.txt\"\r\nContent-Transfer-Encoding: base64\r\nContent-Disposition: attachment; filename=\"file.txt\"\r\nContent-ID: <file.txt>\r\n\r\ndGV4dCBmaWxl",
			testFunc: func(w *partWriter) error {
				return w.writeAttachment(&attachmentPart{
					name: "file.txt",
					copier: func(w io.Writer) error {
						_, _ = w.Write([]byte(`text file`))
						return nil
					},
				})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testWrite(t, tc.testFunc, tc.want)
		})
	}
}

func testWrite(t *testing.T, testFunc func(w *partWriter) error, res string) {
	buf := new(bytes.Buffer)

	w := partWriter{
		w:        buf,
		boundary: "boundary-1",
	}
	err := testFunc(&w)
	require.NoError(t, err)
	assert.Equal(t, res, buf.String())
}
