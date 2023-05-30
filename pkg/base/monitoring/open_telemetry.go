package monitoring

import (
	"context"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/url"
)

type openTelemetry struct {
	tracerProvider trace.TracerProvider
	tracer         trace.Tracer
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(fmt.Sprintf("%s-otel", config.APP_NAME)),
		//semconv.ServiceVersion("0.0.1"), // TODO get current app version: use branch name or commit hash
	)
}

func startOpenTelemetryMonitoring() monitoring {
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

	return &openTelemetry{tracer: tracer}
}

func (m *openTelemetry) startTransaction(ctx context.Context, name string) (interface{}, context.Context) {
	ctx, span := m.tracer.Start(ctx, name)
	return span, ctx
}

func (m *openTelemetry) endTransaction(span interface{}) {
	span.(trace.Span).End()
}

func (m *openTelemetry) setWebRequest(ctx context.Context, transaction interface{}, header http.Header, url *url.URL, method string) {
	panic("not implemented")
}

func (m *openTelemetry) startWebRequest(ctx context.Context, header http.Header, path string, method string) (interface{}, context.Context) {
	attrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String(method),
		// FIXME config attributes
		//semconv.HTTPRequestContentLengthKey.Int(c.Request().Header.ContentLength()),
		//semconv.HTTPSchemeKey.String(utils.CopyString(c.Protocol())),
		//semconv.HTTPTargetKey.String(string(utils.CopyBytes(c.Request().RequestURI()))),
		semconv.HTTPURLKey.String(path),
		////semconv.HTTPUserAgentKey.String(string(utils.CopyBytes(c.Request().Header.UserAgent()))),
		//semconv.NetHostNameKey.String(utils.CopyString(c.Hostname())),
		semconv.NetTransportTCP,
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindServer),
	}
	ctx, span := m.tracer.Start(ctx, fmt.Sprintf("%s %s", method, path), opts...)

	return span, ctx
}

func (m *openTelemetry) setWebResponse(transaction interface{}, w http.ResponseWriter) http.ResponseWriter {
	//TODO implement me
	panic("implement me")
}

func (m *openTelemetry) startTransactionSegment(ctx context.Context, transaction interface{}, name string, attributes map[string]interface{}) interface{} {
	_, span := m.tracer.Start(ctx, name)
	return span
}

func (m *openTelemetry) endTransactionSegment(segment interface{}) {
	segment.(trace.Span).End()
}

func (m *openTelemetry) getTransactionInContext(ctx context.Context) interface{} {
	return trace.SpanFromContext(ctx)
}

func (m *openTelemetry) noticeError(transaction interface{}, err error) {
	transaction.(trace.Span).RecordError(err)
}
