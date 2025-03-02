package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

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
	policyARN := "arn:aws:iam::128114088254:policy/assume-role-policy"

	getPolicyOutput, err := iamClient.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: aws.String(policyARN),
	})
	if err != nil {
		fmt.Printf("failed to create session: %v", err)
	}

	getPolicyVersionOutput, err := iamClient.GetPolicyVersion(&iam.GetPolicyVersionInput{
		PolicyArn: aws.String(policyARN),
		VersionId: getPolicyOutput.Policy.DefaultVersionId,
	})
	if err != nil {
		log.Fatalf("Failed to get policy version: %v", err)
	}

	encodedPolicyDocument := *getPolicyVersionOutput.PolicyVersion.Document
	decodedPolicyDocument, err := url.QueryUnescape(encodedPolicyDocument)
	if err != nil {
		log.Fatalf("Failed to decode policy document: %v", err)
	}

	var policy Policy
	err = json.Unmarshal([]byte(decodedPolicyDocument), &policy)
	if err != nil {
		log.Fatalf("Failed to unmarshal policy: %v", err)
	}

	// Update the policy by adding the new ARN
	for i, statement := range policy.Statement {
		if statement.Action == "sts:AssumeRole" {
			// Check the current type of Resource
			switch resource := statement.Resource.(type) {
			case string:
				// If Resource is a string, convert it to a slice and add the new ARN
				policy.Statement[i].Resource = []string{resource, arnAttachRequest.Arn}
			case []interface{}:
				// If Resource is already a slice, append the new ARN
				resources := make([]string, len(resource))
				for j, r := range resource {
					if str, ok := r.(string); ok {
						resources[j] = str
					}
				}
				resources = append(resources, arnAttachRequest.Arn)
				policy.Statement[i].Resource = resources
			default:
				log.Fatalf("Unexpected Resource type: %T", resource)
			}
		}
	}

	// Marshal the updated policy back to JSON
	updatedPolicyBytes, err := json.Marshal(policy)
	if err != nil {
		log.Fatalf("Failed to marshal updated policy: %v", err)
	}

	// Create a new policy version with the updated document
	_, err = iamClient.CreatePolicyVersion(&iam.CreatePolicyVersionInput{
		PolicyArn:      aws.String(policyARN),
		PolicyDocument: aws.String(strings.ReplaceAll(string(updatedPolicyBytes), " ", "")),
		SetAsDefault:   aws.Bool(true), // Set this version as the default
	})
	if err != nil {
		log.Fatalf("Failed to create policy version: %v", err)
	}

	return util.SuccessResponse(string(updatedPolicyBytes)), nil
}
