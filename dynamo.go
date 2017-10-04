package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
		TableName:              aws.String("BorkedJobs"),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		log.Println(err.Error())
	}
}
