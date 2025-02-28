package main

import (
	"context"
	"net/http"
	"users/handlers"

	"github.com/aws/aws-lambda-go/events"
)

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resouce := req.Resource
	method := req.HTTPMethod

	switch {
	// USERS
	case resouce == "/users" && method == http.MethodPost:
		return handlers.AddUser(ctx, req)
	case resouce == "/users/{userUuid}" && method == http.MethodGet:
		return handlers.GetUser(ctx, req)
	case resouce == "/users/{userUuid}" && method == http.MethodPut:
		return handlers.UpdateUser(ctx, req)
	case resouce == "/users/{userUuid}" && method == http.MethodDelete:
		return handlers.DeleteUser(ctx, req)

	// AUTH
	case resouce == "/sign-in" && method == http.MethodPost:
		return handlers.SignInUser(ctx, req)

	// DEFAULT
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
