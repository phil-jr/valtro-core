package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	// "log"
	// "net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	"companies/util"
)

type ArnAttachRequest struct {
	Arn string `json:"arn"`
}

// Policy represents the overall policy document.
type Policy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// Statement represents a single statement in the policy.
type Statement struct {
	Effect   string      `json:"Effect"`
	Action   string      `json:"Action"`
	Resource interface{} `json:"Resource"`
}

func AttachCompanyRoleArn(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var arnAttachRequest ArnAttachRequest

	json_err := json.Unmarshal([]byte(req.Body), &arnAttachRequest)
	if json_err != nil {
		return util.InputErrorResponse("Invalid JSON"), nil
	}

	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("failed to create session: %v", err)
	}

	iamClient := iam.New(sess)

	roleOutput, err := iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String("alpha-metric-processor-role-616fl245 "),
	})
	if err != nil {
		fmt.Printf("failed to create session: %v", err)
	}

	encodedPolicy := aws.StringValue(roleOutput.Role.AssumeRolePolicyDocument)
	decodedPolicy, err := url.QueryUnescape(encodedPolicy)
	if err != nil {
		fmt.Printf("failed to create session: %v", err)
	}

	return util.SuccessResponse(decodedPolicy), nil
}

func (p *Policy) addResource(newARN string) error {
	for i, stmt := range p.Statement {
		// Check if the Action (or any other indicator) matches what we expect.
		// In this case, we target the sts:AssumeRole statement.
		if stmt.Action == "sts:AssumeRole" {
			switch resource := stmt.Resource.(type) {
			case string:
				// Convert single string to slice
				p.Statement[i].Resource = []string{resource, newARN}
			case []interface{}:
				// Convert []interface{} to []string while appending newARN
				var resources []string
				for _, v := range resource {
					if s, ok := v.(string); ok {
						resources = append(resources, s)
					} else {
						fmt.Printf("unexpected resource type")
					}
				}
				resources = append(resources, newARN)
				p.Statement[i].Resource = resources
			case []string:
				// Directly append if already a slice of strings
				p.Statement[i].Resource = append(resource, newARN)
			default:
				fmt.Printf("unknown type for Resource")
			}
		}
	}
	return nil
}
