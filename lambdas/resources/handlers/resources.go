package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"resources/db"
	"resources/types"
	"time"

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
	endTimestamp := time.Now().UTC()
	startTimestamp := time.Unix(0, 0).UTC()
	companyUuid, err := getMapValue(req.PathParameters, "companyUuid")
	if err != nil {
		return inputErrorResponse(err.Error()), nil
	}

	// THIS IS WHAT SECURES THE ENDPOINT
	if !UserCanAccessEndpoint(req.Headers, companyUuid) {
		return forbiddenError("Missing Authorization header"), nil
	}

	resourceUuid, err := getMapValue(req.PathParameters, "resourceUuid")
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

	rows, err := db.Pool.Query(ctx, db.SelectResouceData, resourceUuid, companyUuid)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
		return internalServerErrorResponse(), nil
	}
	defer rows.Close()

	var metrics []types.Metric
	for rows.Next() {
		var metric types.Metric
		err := rows.Scan(&metric.Name, &metric.Value, &metric.Unit, &metric.Timestamp)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		metrics = append(metrics, metric)
	}

	jsonData, _ := json.MarshalIndent(metrics, "", "  ")

	return successResponseWithBody(string(jsonData)), nil
}

func GetCompanyResourceTotalCost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	totalCost := 0.0
	endTimestamp := time.Now().UTC()
	startTimestamp := time.Unix(0, 0).UTC()

	companyUuid, err := getMapValue(req.PathParameters, "companyUuid")
	if err != nil {
		return inputErrorResponse(err.Error()), nil
	}

	// THIS IS WHAT SECURES THE ENDPOINT
	if !UserCanAccessEndpoint(req.Headers, companyUuid) {
		return forbiddenError("Unauthorized"), nil
	}

	resourceUuid, err := getMapValue(req.PathParameters, "resourceUuid")
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

	rows, err := db.Pool.Query(ctx, db.SelectResouceTotalCost, resourceUuid, companyUuid, startTimestamp, endTimestamp)
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

func GetCompanyResourceCost(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	endTimestamp := time.Now().UTC()
	startTimestamp := time.Unix(0, 0).UTC()

	companyUuid, err := getMapValue(req.PathParameters, "companyUuid")
	if err != nil {
		return inputErrorResponse(err.Error()), nil
	}

	// THIS IS WHAT SECURES THE ENDPOINT
	if !UserCanAccessEndpoint(req.Headers, companyUuid) {
		return forbiddenError("Unauthorized"), nil
	}

	resourceUuid, err := getMapValue(req.PathParameters, "resourceUuid")
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

	var costMetrics []types.Cost
	for rows.Next() {
		var cost types.Cost
		err := rows.Scan(
			&cost.ResourceID,
			&cost.Cost,
			&cost.Aggregate,
			&cost.StartTimestamp,
			&cost.EndTimestamp,
			&cost.CreatedAt,
		)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}

		costMetrics = append(costMetrics, cost)
	}

	jsonData, _ := json.MarshalIndent(costMetrics, "", "  ")

	return successResponseWithBody(string(jsonData)), nil
}
