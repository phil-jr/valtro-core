package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request received: %s %s %s", req.HTTPMethod, req.Resource, req.QueryStringParameters)

	return Router(ctx, req)
}

func main() {
	lambda.Start(handler)
}
