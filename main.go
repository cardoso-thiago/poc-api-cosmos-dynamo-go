package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/config"
	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/repository"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

func main() {
	cleanup := initTracer()
	defer cleanup(context.Background())
	start := time.Now()
	cfg := config.LoadConfig()

	var repo repository.Repository
	var err error

	switch cfg.CloudProvider {
	case "azure":
		repo, err = repository.NewCosmosRepository(cfg.CosmosEndpoint, cfg.CosmosKey, cfg.CosmosDatabase, cfg.CosmosContainer)
		if err != nil {
			panic(err)
		}
	case "aws":
		repo, err = repository.NewDynamoRepository(cfg.DynamoEndpoint, cfg.DynamoRegion, cfg.DynamoTable)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("CLOUD_PROVIDER inválido. Use 'azure' ou 'aws'.")
		os.Exit(1)
	}

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.GET("/items/:id", func(c *gin.Context) {
		id := c.Param("id")

		ctx, span := otel.Tracer(serviceName).Start(c.Request.Context(), "db.GetItem")
		span.SetAttributes(attribute.String("cloud_provider", cfg.CloudProvider))
		item, err := repo.GetItem(ctx, id)
		defer span.End()

		if err != nil {
			if strings.Contains(err.Error(), "404") || os.IsNotExist(err) {
				c.Status(http.StatusNoContent)
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	elapsed := time.Since(start)
	go func() {
		fmt.Printf("API started in %s\n", elapsed)
	}()

	r.Run(":8888")
}

func initTracer() func(context.Context) error {
	secureOption := otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		panic(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}
