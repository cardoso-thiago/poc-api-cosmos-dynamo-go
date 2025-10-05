package repository

import (
	"context"

	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/models"
)

type Repository interface {
	GetItem(ctx context.Context, id string) (*models.Relationship, error)
}
