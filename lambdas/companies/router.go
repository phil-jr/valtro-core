package main

import (
	"companies/handlers"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resouce := req.Resource
	method := req.HTTPMethod

	switch {
	// COMPANIES
	case resouce == "/companies" && method == http.MethodGet:
		return handlers.GetAllCompanies(ctx, req)

	case resouce == "/companies" && method == http.MethodPost:
		return handlers.AddCompany(ctx, req)

	case resouce == "/companies/{companyUuid}" && method == http.MethodGet:
		return handlers.GetCompany(ctx, req)

	case resouce == "/companies/{companyUuid}" && method == http.MethodDelete:
		return handlers.DeleteCompany(ctx, req)

	case resouce == "/companies/{companyUuid}" && method == http.MethodPut:
		return handlers.UpdateCompany(ctx, req)

	// TEAMS
	case resouce == "/companies/{companyUuid}/teams" && method == http.MethodPost:
		return handlers.AddTeam(ctx, req)

	case resouce == "/companies/{companyUuid}/teams" && method == http.MethodGet:
		return handlers.GetAllCompanyTeams(ctx, req)

	case resouce == "/companies/{companyUuid}/teams/{teamUuid}" && method == http.MethodGet:
		return handlers.GetTeam(ctx, req)

	case resouce == "/companies/{companyUuid}/teams/{teamUuid}" && method == http.MethodDelete:
		return handlers.DeleteTeam(ctx, req)

	case resouce == "/companies/{companyUuid}/teams/{teamUuid}" && method == http.MethodPut:
		return handlers.UpdateTeam(ctx, req)

	// TEAM MEMBERS
	case resouce == "/companies/{companyUuid}/teams/{teamUuid}/members" && method == http.MethodGet:
		return handlers.GetTeamMembers(ctx, req)

	case resouce == "/companies/{companyUuid}/teams/{teamUuid}/members" && method == http.MethodPost:
		return handlers.AddTeamMember(ctx, req)

	case resouce == "/companies/{companyUuid}/teams/{teamUuid}/members" && method == http.MethodDelete:
		return handlers.DeleteTeamMember(ctx, req)

	//ACCOUNTS
	case resouce == "/companies/{companyUuid}/attachRoleArn" && method == http.MethodPost:
		return handlers.AttachCompanyRoleArn(ctx, req)

	//DEFAULT
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error":"Not Found Company Lambda"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}
