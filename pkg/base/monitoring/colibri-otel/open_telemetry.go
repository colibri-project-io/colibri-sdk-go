package colibri_otel

import (
	"context"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring/colibri-monitoring-base"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
)

type MonitoringOpenTelemetry struct {
	tracerProvider trace.TracerProvider
	tracer         trace.Tracer
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(config.APP_NAME),
		semconv.ServiceVersion(config.APP_VERSION),
	)
}

func StartOpenTelemetryMonitoring() colibri_monitoring_base.Monitoring {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		logging.Fatal("Creating OTLP trace exporter: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer(
		"github.com/colibri-project-io/colibri-sdk-go",
		trace.WithInstrumentationVersion(contrib.SemVersion()),
	)

	return &MonitoringOpenTelemetry{tracer: tracer}
}

func (m *MonitoringOpenTelemetry) StartTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	ctx, span := m.tracer.Start(ctx, name)
	return span, ctx
}

func (m *MonitoringOpenTelemetry) EndTransaction(span interface{}) {
	span.(trace.Span).End()
}

func (m *MonitoringOpenTelemetry) StartWebRequest(ctx context.Context, header http.Header, path string, method string) (interface{}, context.Context) {
	attrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String(method),
		semconv.HTTPRequestContentLengthKey.String(header.Get("Content-Length")),
		semconv.HTTPSchemeKey.String(header.Get("X-Protocol")),
		semconv.HTTPTargetKey.String(header.Get("X-Request-URI")),
		semconv.HTTPURLKey.String(path),
		semconv.UserAgentOriginal(header.Get("User-Agent")),
		semconv.NetHostNameKey.String(header.Get("Host")),
		semconv.NetTransportTCP,
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindServer),
	}
	ctx, span := m.tracer.Start(ctx, fmt.Sprintf("%s %s", method, path), opts...)

	return span, ctx
}

func (m *MonitoringOpenTelemetry) StartTransactionSegment(ctx context.Context, name string, attributes map[string]string) interface{} {
	_, span := m.tracer.Start(ctx, name)

	kv := make([]attribute.KeyValue, 0, len(attributes))
	for key, value := range attributes {
		kv = append(kv, attribute.String(key, value))
	}
	span.SetAttributes(kv...)

	return span
}

func (m *MonitoringOpenTelemetry) EndTransactionSegment(segment interface{}) {
	segment.(trace.Span).End()
}

func (m *MonitoringOpenTelemetry) GetTransactionInContext(ctx context.Context) interface{} {
	return trace.SpanFromContext(ctx)
}

func (m *MonitoringOpenTelemetry) NoticeError(transaction interface{}, err error) {
	transaction.(trace.Span).RecordError(err)
	transaction.(trace.Span).SetStatus(codes.Error, err.Error())
}

func (m *MonitoringOpenTelemetry) GetSQLDBDriverName() string {
	driverName, err := otelsql.Register("postgres",
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
		otelsql.WithDatabaseName(os.Getenv(config.SQL_DB_NAME)),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		logging.Fatal("could not get sql db driver name: %v", err)
	}
	return driverName
}
