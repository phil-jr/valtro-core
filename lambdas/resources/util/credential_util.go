package util

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func FetchLambdaClient(roleArn string) *lambda.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error: " + err.Error())
	}

	stsClient := sts.NewFromConfig(cfg)
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("TempSession"),
	}

	result, err := stsClient.AssumeRole(context.TODO(), input)
	if err != nil {
		panic("error assuming role: " + err.Error())
	}

	creds := aws.Credentials{
		AccessKeyID:     *result.Credentials.AccessKeyId,
		SecretAccessKey: *result.Credentials.SecretAccessKey,
		SessionToken:    *result.Credentials.SessionToken,
	}

	assumedRoleCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				creds.AccessKeyID,
				creds.SecretAccessKey,
				creds.SessionToken,
			),
		)),
	)
	if err != nil {
		panic("error assuming role: " + err.Error())
	}

	return lambda.NewFromConfig(assumedRoleCfg)
}
