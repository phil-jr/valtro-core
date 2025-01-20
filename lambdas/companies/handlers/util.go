package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Helper for 500 responses
func internalServerErrorResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"error":"Internal Server Error"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}
