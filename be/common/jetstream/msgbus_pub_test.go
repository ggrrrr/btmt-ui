package jetstream

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	msgbusv1 "github.com/ggrrrr/btmt-ui/be/common/msgbus/v1"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

func Test(t *testing.T) {
	os.Setenv("OTEL_COLLECTOR", "localhost:4317")
	os.Setenv("SERVICE_NAME", "test-service")

	var err error
	ctx := context.Background()
	// wg := sync.WaitGroup{}

	mockTokenGenerotar := token.NewTokenGenerator("test-publisher", token.NewSignerMock())

	err = tracer.ConfigureForTest()
	require.NoError(t, err)
	defer func() {
		tracer.Shutdown(ctx)
		fmt.Println("logger.Shutdown ;)")
	}()

	conn, err := ConnectForTest()
	require.NoError(t, err)
	defer func() {
		conn.conn.Close()
		fmt.Println("conn.conn.Close")
	}()

	id := uuid.New()

	testtPublisher, err := NewCommandPublisher[*msgbusv1.TestCommand](ctx, conn, &msgbusv1.TestCommand{}, mockTokenGenerotar)
	require.NoError(t, err)

	testtPublisher.Publish(ctx, msgbus.Metadata{Id: id}, &msgbusv1.TestCommand{Id: "asd", Name: "asdad"})

}
