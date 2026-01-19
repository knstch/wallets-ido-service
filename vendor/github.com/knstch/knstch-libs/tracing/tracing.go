package tracing

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func InitTracer(serviceName, jaegerHost string) func(ctx context.Context) error {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerHost)),
	)
	if err != nil {
		log.Fatalf("failed to create jaeger exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	tracer = otel.Tracer(serviceName)

	return tp.Shutdown
}

func WithTracing() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			method, _ := ctx.Value(httptransport.ContextKeyRequestMethod).(string)
			path, _ := ctx.Value(httptransport.ContextKeyRequestPath).(string)

			spanName := fmt.Sprintf("%s %s", method, path)
			ctx, span := tracer.Start(ctx, spanName)
			defer span.End()

			return next(ctx, request)
		}
	}
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}
