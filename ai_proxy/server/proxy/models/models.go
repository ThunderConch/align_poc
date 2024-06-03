package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.Must(session.NewSession(&aws.Config{
	Region: aws.String("ap-northeast-2"), // Asia Pacific (Seoul)
})))

type APIKeyInfo struct {
	APIKey string   `json:"APIKey"`
	RPM    int      `json:"RPM"`
	BPM    int64    `json:"BPM"`
	Nodes  []string `json:"Nodes"`
}

func GetAPIKeyInfo(key string) (*APIKeyInfo, error) {
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("APIKeys"),
		Key: map[string]*dynamodb.AttributeValue{
			"APIKey": {
				S: aws.String(key),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var apiKeyInfo APIKeyInfo
	err = dynamodbattribute.UnmarshalMap(result.Item, &apiKeyInfo)
	if err != nil {
		return nil, err
	}

	return &apiKeyInfo, nil
}
