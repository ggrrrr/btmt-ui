package logger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	Config struct {
		Collector string `env:"OTEL_COLLECTOR"`
	}
)

var (
	rootCtx context.Context

	spanKeyPrefix = "btmt"
	otelScopeName = "go.github.com.ggrrrr.btmt-ui"

	tracerProvider trace.TracerProvider

	tracer trace.Tracer

	shutdownFunc func(context.Context) error

	DevConfig = Config{
		Collector: "localhost:4222",
	}
)

func initNoopOtel() {
	// traceProvider = trace.NewNoopTracerProvider()
	Info().Msg("otel.noop.")
	tracerProvider = noop.NewTracerProvider()
	shutdownFunc = func(context.Context) error {
		return nil
	}
	tracer = tracerProvider.Tracer(otelScopeName)
}

func ConfigureOtel(ctx context.Context, appName string, cfg Config) error {

	rootCtx = ctx

	if cfg.Collector == "" {
		initNoopOtel()
		return fmt.Errorf("OTEL_COLLECTOR not set")
	}

	if appName == "" {
		appName = "undefined-app-name"
	}

	serviceNameKey := semconv.ServiceNameKey.String(appName)
	// traceProvider = trace.NewNoopTracerProvider()
	conn, err := grpc.NewClient(cfg.Collector,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		Error(err).Msg("initGrpcOtel.conn")
		initNoopOtel()
		return err
	}

	otelResource, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			serviceNameKey,
		),
	)
	if err != nil {
		Error(err).Msg("initGrpcOtel.otelResource")
		conn.Close()
		initNoopOtel()
		return err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		conn.Close()
		Error(err).Msg("initGrpcOtel.exporter")
		initNoopOtel()
		return err
	}

	Info().Str("addr", cfg.Collector).Msg("otel.grpc.")

	batchProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	sdkTracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(otelResource),
		sdktrace.WithSpanProcessor(batchProcessor),
	)

	otel.SetTracerProvider(sdkTracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer = sdkTracerProvider.Tracer(otelScopeName)
	tracerProvider = sdkTracerProvider
	shutdownFunc = sdkTracerProvider.Shutdown

	return nil
}

func Shutdown() {
	err := shutdownFunc(rootCtx)
	if err != nil {
		Error(err).Msg("otel.Shutdown")
		return
	}
	Info().Msg("otel.Shutdown")
}

func Tracer() trace.Tracer {
	return tracer
}

func TracerProvider() trace.TracerProvider {
	return tracerProvider
}

func Span(ctx context.Context, name string, payload TraceDataExtractor) (context.Context, AppSpan) {
	traceData := TraceDataAppend(TraceDataFromCtx(ctx), payload)
	ctx = context.WithValue(ctx, traceDataCtxKey{}, traceData)

	attributes := SpanAttributes(roles.AuthInfoFromCtx(ctx), traceData)
	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(attributes...))
	return ctx,
		&appSpan{
			span: span,
		}
}

func SpanWithAttributes(ctx context.Context, name string, payload TraceDataExtractor, attribs ...TraceKV) (context.Context, AppSpan) {
	traceData := TraceDataAppend(TraceDataFromCtx(ctx), payload, attribs...)
	ctx = context.WithValue(ctx, traceDataCtxKey{}, traceData)

	attributes := SpanAttributes(roles.AuthInfoFromCtx(ctx), traceData)
	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(attributes...))
	return ctx,
		&appSpan{
			span: span,
		}
}
