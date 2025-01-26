package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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
	err := row.Scan(
		&user.UserId,
		&user.CompanyId,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Admin,
		&user.CreatedAt,
	)
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

	return successResponseWithBody(string(body)), nil
}

func UpdateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func DeleteUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internalServerErrorResponse(), nil
}

func SignInUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var signIn types.SignIn
	json_err := json.Unmarshal([]byte(req.Body), &signIn)
	if json_err != nil {
		return inputErrorResponse("Invalid JSON"), nil
	}

	row := db.Pool.QueryRow(ctx, db.RetrievePasswordHashQuery, signIn.Email)
	err := row.Scan(
		&signIn.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("There was an issue signing in user with: %v", signIn.Email)
			return inputErrorResponseUnauthorized("Incorrect email or password."), nil
		}
		log.Printf("Error executing query: %v", err)
		return internalServerErrorResponse(), nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(signIn.Password), []byte(signIn.Password))
	if err != nil {
		fmt.Println("Password verification failed:", err)
		return inputErrorResponseUnauthorized("Incorrect email or password."), nil
	}

	return successResponse("User sign in successful!"), nil
}
