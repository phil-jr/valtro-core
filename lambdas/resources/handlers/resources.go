package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"resources/db"
	"time"

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
	totalCost := 0.0
	endTimestamp := time.Now().UTC()
	startTimestamp := time.Unix(0, 0).UTC()

	companyUuid, err := getPathParam(req.PathParameters, "companyUuid")
	if err != nil {
		return inputErrorResponse(err.Error()), nil
	}

	resourceUuid, err := getPathParam(req.PathParameters, "resourceUuid")
	if err != nil {
		return inputErrorResponse(err.Error()), nil
	}

	if t, err := parseQueryTime(req.QueryStringParameters, "startTime", startTimestamp); err != nil {
		log.Printf("Error parsing startTime: %v", err)
		return inputErrorResponse("Invalid startTime format"), nil
	} else {
		startTimestamp = t
	}

	if t, err := parseQueryTime(req.QueryStringParameters, "endTime", endTimestamp); err != nil {
		log.Printf("Error parsing endTime: %v", err)
		return inputErrorResponse("Invalid endTime format"), nil
	} else {
		endTimestamp = t
	}

	rows, err := db.Pool.Query(ctx, db.SelectResouceCost, resourceUuid, companyUuid, startTimestamp, endTimestamp)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
		return internalServerErrorResponse(), nil
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
		totalCost += (cost * float64(occurrences))
		fmt.Printf("Cost: %f, Occurrences: %d\n", cost, occurrences)
	}

	payload := map[string]float64{
		"totalCost": totalCost,
	}

	body, err := json.Marshal(payload)
	return successResponseWithBody(string(body)), nil
}
