package api

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stefanprodan/podinfo/pkg/version"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	instrumentationName = "github.com/stefanprodan/podinfo/pkg/api"
)

func (a *Api) InitTracer(ctx context.Context) {
	if viper.GetString("otel-service-name") == "" {
		nop := trace.NewNoopTracerProvider()
		a.Tracer = nop.Tracer(viper.GetString("otel-service-name"))
		return
	}

	client := otlptracegrpc.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		a.Logger.Error("creating OTLP trace exporter", zap.Error(err))
	}

	a.TracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(viper.GetString("otel-service-name")),
			semconv.ServiceVersionKey.String(version.VERSION),
		)),
	)

	otel.SetTracerProvider(a.TracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		b3.New(),
		&jaeger.Jaeger{},
		&ot.OT{},
		&xray.Propagator{},
	))

	a.Tracer = a.TracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(version.VERSION),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
}

func NewOpenTelemetryMiddleware() mux.MiddlewareFunc {
	return otelmux.Middleware(viper.GetString("otel-service-name"))
}
