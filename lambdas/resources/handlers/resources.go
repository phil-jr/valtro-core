package handlers

import (
	"context"
	// "encoding/json"
	// "log"
	// "net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetAllResources(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func GetAllCompanyResources(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func GetCompanyResource(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func UpdateCompanyResource(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func RemoveCompanyResource(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func GetCompanyResourceData(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}
