package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request received WORKING???: %s %s", req.HTTPMethod, req.Resource)

	return Router(
		c,
		tx, req)
}

func main() {
	lambda.Start(handler)
}
