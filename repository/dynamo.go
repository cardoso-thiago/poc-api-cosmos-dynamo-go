package repository

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cardoso-thiago/poc-api-cosmos-dynamo-go/models"
)

type DynamoRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewDynamoRepository(endpoint, region, tableName string) (*DynamoRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
	})
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)
	return &DynamoRepository{db: db, tableName: tableName}, nil
}

func (r *DynamoRepository) GetItem(ctx context.Context, id string) (*models.Relationship, error) {
	result, err := r.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, os.ErrNotExist
	}
	var rel models.Relationship
	b, err := json.Marshal(result.Item)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &rel); err != nil {
		return nil, err
	}
	return &rel, nil
}
