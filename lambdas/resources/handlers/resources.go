package handlers

import (
	"context"
	"fmt"
	"log"
	"resources/db"

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

// TODO: Check if user has permission
func GetCompanyResourceCost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resourceUuid, ok := req.PathParameters["resourceUuid"]
	if !ok {
		return inputErrorResponse("Missing path param!"), nil
	}

	rows, err := db.Pool.Query(ctx, db.SelectResouceCost, resourceUuid)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cost float64
		var occurrences int
		err := rows.Scan(&cost, &occurrences)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		fmt.Printf("Cost: %f, Occurrences: %d\n", cost, occurrences)
	}

	return internalServerErrorResponse(), nil
}
