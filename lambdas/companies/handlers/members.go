package handlers

import (
	"context"
	// "encoding/json"
	// "log"
	// "net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetTeamMembers(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func AddTeamMember(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func DeleteTeamMember(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}
