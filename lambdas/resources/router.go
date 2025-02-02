package main

import (
	"context"
	"log"
	"net/http"
	"resources/handlers"

	"github.com/aws/aws-lambda-go/events"
)

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resouce := req.Resource
	method := req.HTTPMethod

	switch {
	// RESOURCES
	case resouce == "/resources" && method == http.MethodGet:
		return handlers.GetAllResources(ctx, req)

	case resouce == "/companies/{companyUuid}/resources" && method == http.MethodGet:
		return handlers.GetAllCompanyResources(ctx, req)

	case resouce == "/companies/{companyUuid}/resources/{resourceUuid}" && method == http.MethodGet:
		return handlers.GetCompanyResource(ctx, req)

	case resouce == "/companies/{companyUuid}/resources/{resourceUuid}" && method == http.MethodPut:
		return handlers.UpdateCompanyResource(ctx, req)

	case resouce == "/companies/{companyUuid}/resources/{resourceUuid}" && method == http.MethodDelete:
		return handlers.RemoveCompanyResource(ctx, req)

	case resouce == "/companies/{companyUuid}/resources/{resourceUuid}/data" && method == http.MethodGet:
		return handlers.GetCompanyResourceData(ctx, req)

	case resouce == "/companies/{companyUuid}/resources/{resourceUuid}/cost" && method == http.MethodGet:
		return handlers.GetCompanyResourceCost(ctx, req)

	//DEFAULT
	default:
		log.Printf("Incorrect resouce: %v", resouce)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error":"Not Found resources lambda"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}
