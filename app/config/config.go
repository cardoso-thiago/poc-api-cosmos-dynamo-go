package config

import "os"

type Config struct {
	CloudProvider   string
	CosmosEndpoint  string
	CosmosKey       string
	CosmosDatabase  string
	CosmosContainer string
	DynamoEndpoint  string
	DynamoRegion    string
	DynamoTable     string
}

func LoadConfig() *Config {
	return &Config{
		CloudProvider:   os.Getenv("CLOUD_PROVIDER"),
		CosmosEndpoint:  os.Getenv("COSMOS_ENDPOINT"),
		CosmosKey:       os.Getenv("COSMOS_KEY"),
		CosmosDatabase:  os.Getenv("COSMOS_DATABASE"),
		CosmosContainer: os.Getenv("COSMOS_CONTAINER"),
		DynamoEndpoint:  os.Getenv("DYNAMO_ENDPOINT"),
		DynamoRegion:    os.Getenv("DYNAMO_REGION"),
		DynamoTable:     os.Getenv("DYNAMO_TABLE"),
	}
}
