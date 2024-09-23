package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"
)

func initTracer() {
	ctx := context.Background()

	// 创建 OTLP HTTP 导出器
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("alloy:4320"),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	// 创建批处理 Span 处理器
	//bsp := trace.NewBatchSpanProcessor(exporter)

	// 创建 TracerProvider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("order-service"),
		)),
	)

	otel.SetTracerProvider(tp)

	// Set the global propagator to use both TraceContext and Baggage
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}
