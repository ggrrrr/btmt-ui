package tracer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ggrrrr/btmt-ui/be/common/buildversion"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

type (
	TargetCfg struct {
		Addr string `env:"ADDR"`
	}
	ClientCfg struct {
		Target TargetCfg `envPrefix:"TARGET_"`
	}

	Config struct {
		Client ClientCfg `envPrefix:"OTEL_"`
	}

	OTelTracer interface {
		Span(ctx context.Context, name string) (context.Context, OTelSpan)
		SpanWithData(ctx context.Context, name string, e td.TraceDataExtractor) (context.Context, OTelSpan)
		SpanWithAttributes(ctx context.Context, name string, attr ...slog.Attr) (context.Context, OTelSpan)
	}

	OTelSpan interface {
		End(error)
		SetAttributes(attr ...slog.Attr)
	}
)

var (
	cfgLock sync.Mutex = sync.Mutex{}

	otelResource *resource.Resource

	shutDownFunc func(context.Context) error = nil
)

func Configure(ctx context.Context, serviceName string, cfg Config) error {
	cfgLock.Lock()
	defer cfgLock.Unlock()

	fmt.Printf("ltm.tracer: %+v\n", cfg)

	serviceNameAttr := semconv.ServiceNameKey.String(serviceName)
	serviceVersionAttr := semconv.ServiceVersionKey.String(buildversion.BuildVersion())

	out := make([]attribute.KeyValue, 0, 2)
	out = append(out, serviceNameAttr)
	out = append(out, serviceVersionAttr)

	r, err := resource.New(
		ctx,
		resource.WithAttributes(out...),
	)
	if err != nil {
		return err
	}

	otelResource = r

	return connect(cfg)
}

func connect(cfg Config) error {
	ctx := context.Background()

	conn, err := grpc.NewClient(cfg.Client.Target.Addr,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		conn.Close()
		return err
	}

	batchProcessor := sdktrace.NewBatchSpanProcessor(exporter)

	tpo := make([]sdktrace.TracerProviderOption, 0, 3)
	tpo = append(tpo,
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(batchProcessor),
	)
	if otelResource != nil {
		tpo = append(tpo, sdktrace.WithResource(otelResource))
	}

	sdkTracerProvider := sdktrace.NewTracerProvider(tpo...)

	shutDownFunc = sdkTracerProvider.Shutdown

	otel.SetTracerProvider(sdkTracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

func Shutdown(ctx context.Context) error {
	if shutDownFunc == nil {
		return nil
	}
	return shutDownFunc(ctx)
}

func Tracer(scopeName string) OTelTracer {
	if shutDownFunc == nil {
		return &noopTracer{}
	}

	return &realTracer{
		tracer: otel.Tracer(scopeName),
	}
}
