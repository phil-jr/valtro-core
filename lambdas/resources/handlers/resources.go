package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"resources/db"
	"resources/types"
	"resources/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/jackc/pgx/v5"
)

func GetAllResources(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return util.InternalServerErrorResponse(), nil
}

func GetAllCompanyResources(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	companyUuid, err := util.GetMapValue(req.PathParameters, "companyUuid")
	if err != nil {
		return util.InputErrorResponse(err.Error()), nil
	}

	// THIS IS WHAT SECURES THE ENDPOINT
	if !util.UserCanAccessEndpoint(req.Headers, companyUuid) {
		return util.ForbiddenError("Missing Authorization header"), nil
	}

	rows, err := db.Pool.Query(ctx, db.SelectAllCompanyResources, companyUuid)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
		return util.InternalServerErrorResponse(), nil
	}
	defer rows.Close()

	var resources []types.Resource
	for rows.Next() {
		var resource types.Resource
		err := rows.Scan(&resource.ResourceID, &resource.ResourceName, &resource.ResourceType, &resource.CreatedAt)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
			continue // Update with real error
		}
		resources = append(resources, resource)
	}

	jsonData, _ := json.MarshalIndent(resources, "", "  ")

	return util.SuccessResponseWithBody(string(jsonData)), nil
}

func GetCompanyResource(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return util.InternalServerErrorResponse(), nil
}

func UpdateCompanyResource(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return util.InternalServerErrorResponse(), nil
}

func RemoveCompanyResource(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return util.InternalServerErrorResponse(), nil
}

func GetCompanyResourceData(ctx context.Context, req events.APIGatewayProxyRequest, companyUuid string) (events.APIGatewayProxyResponse, error) {
	endTimestamp := time.Now().UTC()
	startTimestamp := time.Unix(0, 0).UTC()

	aggregate := 5
	if aggStr, err := util.GetMapValue(req.QueryStringParameters, "aggregate"); err == nil {
		if aggInt, err := strconv.Atoi(aggStr); err == nil {
			aggregate = aggInt
		}
	}

	resourceUuid, err := util.GetMapValue(req.PathParameters, "resourceUuid")
	if err != nil {
		return util.InputErrorResponse(err.Error()), nil
	}

	if t, err := util.ParseQueryTime(req.QueryStringParameters, "startTime", startTimestamp); err != nil {
		log.Printf("Error parsing startTime: %v", err)
		return util.InputErrorResponse("Invalid startTime format"), nil
	} else {
		startTimestamp = t
	}

	if t, err := util.ParseQueryTime(req.QueryStringParameters, "endTime", endTimestamp); err != nil {
		log.Printf("Error parsing endTime: %v", err)
		return util.InputErrorResponse("Invalid endTime format"), nil
	} else {
		endTimestamp = t
	}

	log.Printf("Start Timestamp: %v | End Timestamp: %v", startTimestamp, endTimestamp)
	rows, err := db.Pool.Query(ctx, db.SelectResouceData, resourceUuid, companyUuid)
	ct := 0
	if err != nil {
		log.Fatalf("Query failed: %v", err)
		return util.InternalServerErrorResponse(), nil
	}
	defer rows.Close()

	var metrics []types.Metric
	for rows.Next() {
		var metric types.Metric
		err := rows.Scan(&metric.Name, &metric.Value, &metric.Unit, &metric.Aggregate, &metric.Timestamp)
		if err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		metrics = append(metrics, metric)
		ct = ct + 1
	}
	log.Printf("rownum: %v", len(metrics))

	evenMetrcis, err := util.EvenlyBucketMetrics(metrics, aggregate)
	if err != nil {
		log.Printf("Even bucket fail: %v", err)
		return util.InputErrorResponse(err.Error()), nil
	}

	jsonData, _ := json.MarshalIndent(evenMetrcis, "", "  ")

	return util.SuccessResponseWithBody(string(jsonData)), nil
}

func GetCompanyResourceTotalCost(ctx context.Context, req events.APIGatewayProxyRequest, companyUuid string) (events.APIGatewayProxyResponse, error) {
	totalCost := 0.0
	endTimestamp := time.Now().UTC()
	startTimestamp := time.Unix(0, 0).UTC()

	resourceUuid, err := util.GetMapValue(req.PathParameters, "resourceUuid")
	if err != nil {
		return util.InputErrorResponse(err.Error()), nil
	}

	if t, err := util.ParseQueryTime(req.QueryStringParameters, "startTime", startTimestamp); err != nil {
		log.Printf("Error parsing startTime: %v", err)
		return util.InputErrorResponse("Invalid startTime format"), nil
	} else {
		startTimestamp = t
	}

	if t, err := util.ParseQueryTime(req.QueryStringParameters, "endTime", endTimestamp); err != nil {
		log.Printf("Error parsing endTime: %v", err)
		return util.InputErrorResponse("Invalid endTime format"), nil
	} else {
		endTimestamp = t
	}

	rows, err := db.Pool.Query(ctx, db.SelectResouceTotalCost, resourceUuid, companyUuid, startTimestamp, endTimestamp)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
		return util.InternalServerErrorResponse(), nil
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
	return util.SuccessResponseWithBody(string(body)), nil
}

func GetCompanyResourceCost(ctx context.Context, req events.APIGatewayProxyRequest, companyUuid string) (events.APIGatewayProxyResponse, error) {
	endTimestamp := time.Now().UTC()
	startTimestamp := time.Unix(0, 0).UTC()

	resourceUuid, err := util.GetMapValue(req.PathParameters, "resourceUuid")
	if err != nil {
		return util.InputErrorResponse(err.Error()), nil
	}

	if t, err := util.ParseQueryTime(req.QueryStringParameters, "startTime", startTimestamp); err != nil {
		log.Printf("Error parsing startTime: %v", err)
		return util.InputErrorResponse("Invalid startTime format"), nil
	} else {
		startTimestamp = t
	}

	if t, err := util.ParseQueryTime(req.QueryStringParameters, "endTime", endTimestamp); err != nil {
		log.Printf("Error parsing endTime: %v", err)
		return util.InputErrorResponse("Invalid endTime format"), nil
	} else {
		endTimestamp = t
	}

	rows, err := db.Pool.Query(ctx, db.SelectResouceCost, resourceUuid, companyUuid, startTimestamp, endTimestamp)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
		return util.InternalServerErrorResponse(), nil
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

	return util.SuccessResponseWithBody(string(jsonData)), nil
}

func UpdateCompanyResourceInfra(ctx context.Context, req events.APIGatewayProxyRequest, companyUuid string) (events.APIGatewayProxyResponse, error) {

	var resourceConfigs types.ResourceConfigs
	err := json.Unmarshal([]byte(req.Body), &resourceConfigs)
	if err != nil {
		return util.InputErrorResponse("Invalid JSON"), err
	}

	resourceUuid, err := util.GetMapValue(req.PathParameters, "resourceUuid")
	if err != nil {
		return util.InputErrorResponse("Invalid URL"), err
	}

	var resource types.ResourceWithArn
	row := db.Pool.QueryRow(
		ctx,
		db.SelectCompanyResource,
		companyUuid,
		resourceUuid,
	)

	err = row.Scan(&resource.ResourceID, &resource.ResourceName, &resource.RoleArn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("Resource not found in DB: company=%s, resource=%s", companyUuid, resourceUuid)
			return util.NotFoundErrorResponse("Resource not found."), nil // Return nil error for 404
		}
		// For any other error, it's internal
		log.Printf("Error scanning resource from DB: %v", err)
		return util.InternalServerErrorResponse(), err // Return 5xx and the actual error
	}

	lambdaClient := util.FetchLambdaClient(resource.RoleArn)
	getInput := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(resource.ResourceName),
	}

	log.Printf("Attempting to get configuration for function: %s\n", resource.ResourceName)

	// Call the GetFunctionConfiguration API
	currentConfigs, err := lambdaClient.GetFunctionConfiguration(ctx, getInput)
	if err != nil {
		log.Printf("Error getting function configuration for %s: %v", resource.ResourceName, err)
		return util.InternalServerErrorResponse(), err
	}
	log.Printf("Current Config - Memory: %v, Timeout: %v", currentConfigs.MemorySize, currentConfigs.Timeout)

	updateNeeded := false
	updateInput := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(resource.ResourceName),
	}

	// Check MemoryMB
	if resourceConfigs.MemoryMB != nil {
		currentMemory := int32(0)
		if currentConfigs.MemorySize != nil {
			currentMemory = *currentConfigs.MemorySize
		}
		if *resourceConfigs.MemoryMB != currentMemory {
			updateInput.MemorySize = aws.Int32(*resourceConfigs.MemoryMB)
			updateNeeded = true
			log.Printf("Planning update for MemorySize from %d to: %d\n", currentMemory, *resourceConfigs.MemoryMB)
		}
	}

	// Check Timeout
	if resourceConfigs.Timeout != nil {
		currentTimeout := int32(0)
		if currentConfigs.Timeout != nil {
			currentTimeout = *currentConfigs.Timeout
		}
		if *resourceConfigs.Timeout != currentTimeout {
			updateInput.Timeout = aws.Int32(*resourceConfigs.Timeout)
			updateNeeded = true
			log.Printf("Planning update for Timeout from %d to: %d\n", currentTimeout, *resourceConfigs.Timeout)
		}
	}

	// --- Execute Update (only if needed) ---
	if updateNeeded { // <<< Check the flag here
		log.Printf("Attempting to update configuration for function: %s\n", resource.ResourceName)
		updateResult, err := lambdaClient.UpdateFunctionConfiguration(ctx, updateInput) // Use ctx
		if err != nil {
			log.Printf("Error updating function configuration: %v", err)
			return util.InternalServerErrorResponse(), err
		}
		log.Printf("Update submitted. Function: %s, Status: %s\n", aws.ToString(updateResult.FunctionName), updateResult.State)
		return util.SuccessResponse("Lambda configuration updated successfully."), nil
	} else {
		log.Printf("No configuration changes needed for function: %s based on provided input.\n", resource.ResourceName)
		return util.SuccessResponse("No configuration changes needed."), nil
	}

}
