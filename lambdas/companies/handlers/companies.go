package handlers

import (
	"context"
	// "encoding/json"
	// "log"
	// "net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetAllCompanies(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func AddCompany(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func GetCompany(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func DeleteCompany(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func UpdateCompany(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}
