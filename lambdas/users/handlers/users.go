package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"users/db"
	"users/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

func AddUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var user types.User
	err := json.Unmarshal([]byte(req.Body), &user)
	user.UserId = uuid.New().String()
	user.CreatedAt = time.Now()
	if err != nil {
		return inputErrorResponse("Invalid JSON"), nil
	}

	var newID string
	err = db.Pool.QueryRow(
		ctx,
		db.InsertUserQuery,
		user.UserId,
		user.CompanyId,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Email,
		user.Admin,
		user.CreatedAt,
	).Scan(&newID)
	if err != nil {
		log.Printf("Failed to insert record and get new ID: %v\n", err)
		return internalServerErrorResponse(), nil
	}

	return successResponse("User created!"), nil
}

func GetUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func UpdateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func DeleteUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}
