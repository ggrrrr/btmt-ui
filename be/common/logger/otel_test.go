package logger

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

var testServiceName = semconv.ServiceNameKey.String("test-service")

type person struct {
	Id   int
	Name string
}

func (d person) Extract() TraceData {
	return map[string]TraceValue{
		"person.id":   TraceValueInt(d.Id),
		"person.name": TraceValueString(d.Name),
	}
}

type testError struct {
	id int
}

func (t testError) Error() string {
	return fmt.Sprintf("error id %v", t.id)
}

func smallerWork(ctx context.Context, payload person) error {
	ii := 0
	var isOdd *int = &ii
	var err error
	_, span := Span(ctx, "smallerWork111", payload)

	defer func() {
		fmt.Printf("%v isOdd error2 %v\n", payload, *isOdd)
	}()
	defer fmt.Printf("%v isOdd error1 %v\n", payload, *isOdd)

	defer span.End(err)

	if payload.Id%2 > 0 {
		asd := 1
		*isOdd = asd
		err = testError{id: payload.Id}
	}

	return err
}

func smallWork(ctx context.Context, payload person) {
	var err error

	ctx, span := Span(ctx, "very.smallWork", payload)
	defer func() {
		span.End(err)
	}()
	// defer span.End(err)

	err = smallerWork(ctx, payload)

}

func TestOtel(t *testing.T) {
	var err error

	os.Setenv("OTEL_COLLECTOR", "localhost:4317")
	os.Setenv("SERVICE_NAME", "test-service")

	rootCtx := context.Background()

	authInfo := roles.AuthInfo{
		Subject: "subject",
		Realm:   "btmt.io",
		Device: roles.Device{
			RemoteAddr: "localhost",
			DeviceInfo: "curl",
		},
	}

	ctx := roles.CtxWithAuthInfo(rootCtx, authInfo)

	err = ConfigureOtel(rootCtx)
	require.NoError(t, err)

	ctx, span := Span(ctx, "main.Method", nil)

	// defer func() {
	// 	span.End(err)
	// }()

	for i := 0; i < 10; i++ {
		p := person{Id: i, Name: fmt.Sprintf("asd %v", i)}
		smallWork(ctx, p)
	}
	span.End(err)

	time.Sleep(1 * time.Second)

	Shutdown()

}
