package main

import (
	"context"
	"log"
	"net/http"
	"users/handlers"

	"github.com/aws/aws-lambda-go/events"
)

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.Path
	method := req.HTTPMethod

	log.Printf("path: %s, method: %s, pathParams: %v", path, method, req.PathParameters)

	switch {
	// USERS
	case path == "/users" && method == http.MethodPost:
		return handlers.AddUser(ctx, req)
	case path == "/users/{userUuid}" && method == http.MethodGet:
		return handlers.GetUser(ctx, req)
	case path == "/users/{userUuid}" && method == http.MethodPut:
		return handlers.UpdateUser(ctx, req)
	case path == "/users/{userUuid}" && method == http.MethodDelete:
		return handlers.DeleteUser(ctx, req)

	// AUTH
	case path == "/sign-in" && method == http.MethodDelete:
		return handlers.DeleteUser(ctx, req)

	//DEFAULT
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error":"Not Found auth lambda"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}
