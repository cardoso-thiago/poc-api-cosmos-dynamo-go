package repository

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/models"
)

type CosmosRepository struct {
	client    *azcosmos.Client
	database  string
	container string
}

func NewCosmosRepository(endpoint, key, database, container string) (*CosmosRepository, error) {
	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		return nil, err
	}

	// Cria um http.Client customizado que ignora a validação de certificado
	// Atencão, deve ser usado apenas para ambiente de testes
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Usa o client HTTP customizado via azcore.ClientOptions
	client, err := azcosmos.NewClientWithKey(endpoint, cred, &azcosmos.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: httpClient,
		},
	})

	if err != nil {
		return nil, err
	}

	return &CosmosRepository{
		client:    client,
		database:  database,
		container: container,
	}, nil
}

func (r *CosmosRepository) GetItem(ctx context.Context, id string) (*models.Relationship, error) {
	cont, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return nil, err
	}

	pk := azcosmos.NewPartitionKeyString(id)
	item, err := cont.ReadItem(ctx, pk, id, nil)
	if err != nil {
		return nil, err
	}

	var rel models.Relationship
	if err := json.Unmarshal(item.Value, &rel); err != nil {
		return nil, err
	}
	return &rel, nil
}
