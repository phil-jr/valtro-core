package main

import (
	"companies/handlers"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.Path
	method := req.HTTPMethod

	switch {
	// COMPANIES
	case path == "/companies" && method == http.MethodGet:
		return handlers.GetAllCompanies(ctx, req)

	case path == "/companies" && method == http.MethodPost:
		return handlers.AddCompany(ctx, req)

	case path == "/companies/{companyUuid}" && method == http.MethodGet:
		return handlers.GetCompany(ctx, req)

	case path == "/companies/{companyUuid}" && method == http.MethodDelete:
		return handlers.DeleteCompany(ctx, req)

	case path == "/companies/{companyUuid}" && method == http.MethodPut:
		return handlers.UpdateCompany(ctx, req)

	// TEAMS
	case path == "/companies/{companyUuid}/teams" && method == http.MethodPost:
		return handlers.AddTeam(ctx, req)

	case path == "/companies/{companyUuid}/teams" && method == http.MethodGet:
		return handlers.GetAllCompanyTeams(ctx, req)

	case path == "/companies/{companyUuid}/teams/{teamUuid}" && method == http.MethodGet:
		return handlers.GetTeam(ctx, req)

	case path == "/companies/{companyUuid}/teams/{teamUuid}" && method == http.MethodDelete:
		return handlers.DeleteTeam(ctx, req)

	case path == "/companies/{companyUuid}/teams/{teamUuid}" && method == http.MethodPut:
		return handlers.UpdateTeam(ctx, req)

	// TEAM MEMBERS
	case path == "/companies/{companyUuid}/teams/{teamUuid}/members" && method == http.MethodGet:
		return handlers.GetTeamMembers(ctx, req)

	case path == "/companies/{companyUuid}/teams/{teamUuid}/members" && method == http.MethodPost:
		return handlers.AddTeamMember(ctx, req)

	case path == "/companies/{companyUuid}/teams/{teamUuid}/members" && method == http.MethodDelete:
		return handlers.DeleteTeamMember(ctx, req)

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
