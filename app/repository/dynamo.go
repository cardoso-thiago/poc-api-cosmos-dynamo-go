package repository

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/models"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type DynamoRepository struct {
	db        *dynamodb.Client
	tableName string
}

func NewDynamoRepository(endpoint, region, tableName string) (*DynamoRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		return nil, err
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &DynamoRepository{db: client, tableName: tableName}, nil
}

func (r *DynamoRepository) GetItem(ctx context.Context, id string) (*models.Relationship, error) {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, os.ErrNotExist
	}

	var rel models.Relationship
	if err := attributevalue.UnmarshalMap(out.Item, &rel); err != nil {
		return nil, err
	}

	return &rel, nil
}
