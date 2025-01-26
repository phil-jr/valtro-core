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
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func AddUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var user types.User
	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		return inputErrorResponse("Invalid JSON"), nil
	}
	user.UserId = uuid.New().String()
	user.CreatedAt = time.Now()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hash)

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
	userUuid, ok := req.PathParameters["userUuid"]
	var user types.User
	if !ok {

		return inputErrorResponse("Missing path param!"), nil
	}

	row := db.Pool.QueryRow(ctx, db.SelectUserQuery, userUuid)
	err := row.Scan(&user.UserId, &user.CompanyId, &user.FirstName, &user.LastName, &user.Password, &user.Email, &user.Admin, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No user found with ID: %v", userUuid)
			return inputErrorResponse("User not found."), nil
		}
		log.Printf("Error executing query: %v", err)
		return internalServerErrorResponse(), nil
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user: %v", err)
		return internalServerErrorResponse(), nil
	}

	return successResponse(string(body)), nil
}

func UpdateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func DeleteUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}
