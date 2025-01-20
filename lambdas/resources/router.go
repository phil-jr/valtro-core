package main

import (
	"context"
	"net/http"
	"resources/handlers"

	"github.com/aws/aws-lambda-go/events"
)

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.Path
	method := req.HTTPMethod

	switch {
	// RESOURCES
	case path == "/resources" && method == http.MethodGet:
		return handlers.GetAllResources(ctx, req)

	case path == "/companies/{companyUuid}/resources" && method == http.MethodGet:
		return handlers.GetAllCompanyResources(ctx, req)

	case path == "/companies/{companyUuid}/resources/{resourceUuid}" && method == http.MethodGet:
		return handlers.GetCompanyResource(ctx, req)

	case path == "/companies/{companyUuid}/resources/{resourceUuid}" && method == http.MethodPut:
		return handlers.UpdateCompanyResource(ctx, req)

	case path == "/companies/{companyUuid}/resources/{resourceUuid}" && method == http.MethodDelete:
		return handlers.RemoveCompanyResource(ctx, req)

	case path == "/companies/{companyUuid}/resources/{resourceUuid}/data" && method == http.MethodPut:
		return handlers.GetCompanyResourceData(ctx, req)

	//DEFAULT
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error":"Not Found resources lambda"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}
