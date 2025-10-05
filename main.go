package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/config"
	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/repository"
	"github.com/gin-gonic/gin"
)

func main() {
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

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.GET("/items/:id", func(c *gin.Context) {
		id := c.Param("id")
		item, err := repo.GetItem(context.Background(), id)
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
