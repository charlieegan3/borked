package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func getCount(key string) (int, error) {
	svc := dynamodb.New(session.New())

	query := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Count": {
				S: aws.String(key),
			},
		},
		TableName: aws.String("borked-counts"),
	}

	result, err := svc.GetItem(query)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(*result.Item["Value"].N)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func setCount(key string, value int) error {
	svc := dynamodb.New(session.New())
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#V": aws.String("Value"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				N: aws.String(fmt.Sprintf("%v", value)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"Count": {
				S: aws.String(key),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String("borked-counts"),
		UpdateExpression: aws.String("SET #V = :v"),
	}

	_, err := svc.UpdateItem(input)
	return err
}

// LogJob saves the results of a job to DynamoDB
func LogJob(rootURL url.URL, linkCount int, source string, userAgent string) {
	// temp workaround to avoid running this in test mode
	if flag.Lookup("test.v") != nil {
		return
	}

	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%v,%v", rootURL.String(), time.Now().UnixNano())))
	jobID := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	svc := dynamodb.New(session.New())
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"JobId": {
				S: aws.String(jobID),
			},
			"RootURL": {
				S: aws.String(rootURL.String()),
			},
			"Count": {
				N: aws.String(fmt.Sprintf("%v", linkCount)),
			},
			"Source": {
				S: aws.String(source),
			},
			"UserAgent": {
				S: aws.String(userAgent),
			},
			"CreatedAt": {
				S: aws.String(time.Now().String()),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("borked-jobs"),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	completedCount, err := getCount("Completed")
	requestsCount, err := getCount("Requests")
	if err != nil {
		log.Println(err.Error())
		return
	}

	setCount("Completed", completedCount+linkCount)
	setCount("Requests", requestsCount+1)
}
