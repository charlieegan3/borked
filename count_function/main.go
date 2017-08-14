package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
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

// Handle is the exported handler called by AWS Lambda.
var Handle apigatewayproxy.Handler

func countHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://borked.charlieegan3.com")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	completedCount, err := getCount("Completed")
	requestsCount, err := getCount("Requests")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseData := struct {
		Completed int `json:"completed"`
		Requests  int `json:"requests"`
	}{
		completedCount,
		requestsCount,
	}

	jsonResult, _ := json.Marshal(responseData)
	w.Write(jsonResult)
}

func init() {
	ln := net.Listen()

	// Amazon API Gateway binary media types are supported out of the box.
	// If you don't send or receive binary data, you can safely set it to nil.
	Handle = apigatewayproxy.New(ln, nil).Handle

	// Any Go framework complying with the Go http.Handler interface can be used.
	// This includes, but is not limited to, Vanilla Go, Gin, Echo, Gorrila, Goa, etc.
	go http.Serve(ln, http.HandlerFunc(countHandler))
}
