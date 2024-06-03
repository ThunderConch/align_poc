package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"ai-proxy/key"
)

var db = dynamodb.New(session.Must(session.NewSession(&aws.Config{
	Region: aws.String("ap-northeast-2"),
})))

func GetAPIKey(key string) (string, error) {
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("APIKeys"),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})

	if err != nil {
		return "", err
	}

	return *result.Item["Value"].S, nil
}

func GetAPIKeyInfo(k string) (*key.KeyInfo, error) {
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("APIKeys"),
		Key: map[string]*dynamodb.AttributeValue{
			"APIKey": {
				S: aws.String(k),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var apiKeyInfo key.KeyInfo
	err = dynamodbattribute.UnmarshalMap(result.Item, &apiKeyInfo)
	if err != nil {
		return nil, err
	}

	return &apiKeyInfo, nil
}
