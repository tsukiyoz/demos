package otel

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func TestServer(t *testing.T) {
	res, err := newResource("demo", "v0.0.1")
	require.NoError(t, err)

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	traceProvider, err := newTraceProvider(res)
	require.NoError(t, err)
	defer traceProvider.Shutdown(context.Background())
	otel.SetTracerProvider(traceProvider)

	server := gin.Default()
	server.GET("/test", func(gCtx *gin.Context) {
		tracer := otel.Tracer("opentelemetry")
		var ctx context.Context = gCtx
		ctx, span := tracer.Start(ctx, "top-span")
		defer span.End()
		span.AddEvent("event-1")
		time.Sleep(time.Second)
		ctx, subSpan := tracer.Start(ctx, "sub-span")
		defer subSpan.End()
		time.Sleep(time.Millisecond * 300)
		subSpan.SetAttributes(attribute.String("key1", "value1"))
		gCtx.String(http.StatusOK, "OK")
	})

	server.Run(":8082")
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
	))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
	)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res))

	return traceProvider, nil
}
