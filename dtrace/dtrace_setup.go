package main

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	COLLECTOR_URL = "aternity-perf-awplatform.aw.k8sw.dev.activenetwork.com"
)

type compositeExporter struct {
	exporters []trace.SpanExporter
}

func (e *compositeExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for _, exporter := range e.exporters {
		if err := exporter.ExportSpans(ctx, spans); err != nil {
			return err
		}
	}

	return nil
}

func (e *compositeExporter) Shutdown(ctx context.Context) error {
	for _, exporter := range e.exporters {
		if err := exporter.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

func serviceName(subName string) string {
	return "active.tax-service." + subName
}

func setupTraceprovider(subName string) {
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName(subName)),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		panic(err)
	}

	otPropagator := ot.OT{}
	otel.SetTextMapPropagator(otPropagator)

	opts := []otlptracehttp.Option{otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(COLLECTOR_URL)}
	client := otlptracehttp.NewClient(opts...)
	httpExporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(err)
	}
	stdoutExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	cExporter := compositeExporter{
		exporters: []trace.SpanExporter{stdoutExporter, httpExporter},
	}

	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(&cExporter)),
		trace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)
}
