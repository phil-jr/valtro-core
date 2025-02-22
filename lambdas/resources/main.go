package main

import (
	"context"
	"log"
	"resources/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request received: %s %s %s", req.HTTPMethod, req.Resource, req.QueryStringParameters)

	// If an endpoint doesn't require a companyUuid, do something
	companyUuid, err := util.GetMapValue(req.PathParameters, "companyUuid")
	if err != nil {
		return util.InputErrorResponse(err.Error()), nil
	}

	// THIS IS WHAT SECURES THE ENDPOINT
	if !util.UserCanAccessEndpoint(req.Headers, companyUuid) {
		return util.ForbiddenError("Missing Authorization header"), nil
	}

	return Router(ctx, req, companyUuid)
}

func main() {
	lambda.Start(handler)
}
