package ltm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type TestWriteCloser struct {
	bytes []byte
	buf   *bytes.Buffer
}

var _ (io.WriteCloser) = (*TestWriteCloser)(nil)

func (twc *TestWriteCloser) Write(p []byte) (n int, err error) {
	return twc.buf.Write(p)
}

func (twc *TestWriteCloser) Close() error {
	return nil
}

func NewTestWriter() *TestWriteCloser {
	out := &TestWriteCloser{
		bytes: []byte{},
	}
	out.buf = bytes.NewBuffer(out.bytes)

	return out
}

func TestLog(t *testing.T) {
	os.Setenv("TEST_LOG_LEVEL", "INFO")
	// os.Setenv("TEST_LOG_FORMAT", "JSON")

	authInfo := roles.CreateSystemAdminUser("local", "me", app.Device{RemoteAddr: "addr"})
	ctx := roles.CtxWithAuthInfo(context.Background(), authInfo)

	// ctx := tracedata.

	wc := NewTestWriter()
	defer func() {
		fmt.Printf("REC: %v", wc.buf)
	}()

	log := newLogger("test", wc)
	// log := newLogger("test", os.Stderr)

	log.DebugCtx(ctx).Msg("asd")

	// log.loggger.

}
