package handlers

import (
	"context"
	"encoding/json"

	// "log"
	// "net/http"

	"github.com/aws/aws-lambda-go/events"

	"companies/util"
)

type ArnAttachRequest struct {
	Arn string `json:"arn"`
}

func AttachCompanyRoleArn(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var arnAttachRequest ArnAttachRequest

	json_err := json.Unmarshal([]byte(req.Body), &arnAttachRequest)
	if json_err != nil {
		return util.InputErrorResponse("Invalid JSON"), nil
	}
	return util.SuccessResponse(arnAttachRequest.Arn), nil
}
