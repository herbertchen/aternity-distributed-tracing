package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	PORT = ":8080"
	ROLE = "server"
)

func main() {
	setupTraceprovider(ROLE)

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName(ROLE)))
	r.GET("/dtrace", func(c *gin.Context) {
		_, tracerOne := otel.Tracer(serviceName("dtrace-server")).Start(c.Request.Context(), "span1",
			trace.WithAttributes(attribute.String("key1", "value1")))
		defer tracerOne.End()

		time.Sleep(1 * time.Second)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(PORT) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
