package trace

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"io"
)

var tracer = trace.NewNoopTracerProvider().Tracer("non initialized tracer")

func Get() trace.Tracer {
	return tracer
}

type closer struct {
	provider *tracesdk.TracerProvider
	exporter tracesdk.SpanExporter
}

func (c closer) Close() error {
	if err := c.provider.ForceFlush(context.TODO()); err != nil {
		return fmt.Errorf("flush provider: %w", err)
	}
	if err := c.provider.Shutdown(context.TODO()); err != nil {
		return fmt.Errorf("shutdown provider: %w", err)
	}
	if err := c.exporter.Shutdown(context.TODO()); err != nil {
		return fmt.Errorf("shutdown exporter: %w", err)
	}
	return nil
}

func InitTracer(logger logr.Logger, jaegerURL string, serviceName string) (io.Closer, error) {
	otel.SetLogger(logger)

	propagators := propagation.NewCompositeTextMapPropagator(b3.New(), propagation.TraceContext{})
	otel.SetTextMapPropagator(propagators)

	exporter, err := NewJaegerExporter(jaegerURL)
	if err != nil {
		return nil, fmt.Errorf("initialize exporter: %w", err)
	}

	tp, err := NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("initialize provider: %w", err)
	}
	otel.SetTracerProvider(tp)

	tracer = tp.Tracer(serviceName)

	// TODO set tracing rate
	return closer{
		provider: tp,
		exporter: exporter,
	}, nil
}

// NewJaegerExporter creates new jaeger exporter
//
//	url example - http://localhost:14268/api/traces
func NewJaegerExporter(url string) (tracesdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}

func NewTraceProvider(exp tracesdk.SpanExporter, ServiceName string) (*tracesdk.TracerProvider, error) {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	), nil
}
