package handlers

import (
	"fmt"
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

// Helper for 400 responses
func inputErrorResponse(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", msg),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

// Helper for 401 Unauthorized responses
func inputErrorResponseUnauthorized(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", msg),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

// Helper for 403 Forbidden responses
func inputErrorResponseForbidden(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusForbidden,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", msg),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

// Helper for 200 responses
func successResponse(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", msg),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

// Wrapper for 200 responces with a body
func successResponseWithBody(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}
