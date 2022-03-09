package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

const (
	CLIENT_ROLE = "client"
)

func TestDTrace(t *testing.T) {
	setupTraceprovider(CLIENT_ROLE)
	client := resty.New()
	client.SetDebug(true)

	otPropagator := ot.OT{}
	for i := 0; i < 2; i++ {
		dTraceR := client.R()
		spanConext, span := otel.Tracer("tax-rate").Start(context.Background(), serviceName("dtrace-client"))
		span.SetAttributes(attribute.String("key1", "value1"),
			attribute.Bool("boolkey", true),
			attribute.IntSlice("intarraykey", []int{1, 2, 3}))
		otPropagator.Inject(spanConext, propagation.HeaderCarrier(dTraceR.Header))
		defer span.End()
		rsp, err := dTraceR.Get(fmt.Sprintf("http://localhost%s/dtrace", PORT))
		if err != nil {
			panic(err)
		}
		log.Println(rsp.StatusCode())
	}
}
