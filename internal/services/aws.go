package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

type AWSService struct {
	Client *cloudtrail.Client
}

func NewAWSService(accessKey, secretKey string) (*AWSService, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),

		config.WithCredentialsProvider(aws.NewCredentialsCache(
			aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     accessKey,
					SecretAccessKey: secretKey,
					SessionToken:    "",
					Source:          "HardcodedCredentials",
				}, nil
			}),
		)),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == cloudtrail.ServiceID {
					return aws.Endpoint{URL: "http://localhost:4566"}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown service: %s", service)
			},
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}
	client := cloudtrail.NewFromConfig(cfg)
	return &AWSService{Client: client}, nil
}

func (a *AWSService) FetchLogs(ctx context.Context, startTimeStr, endTimeStr string) ([]types.Event, error) {
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid startTime format: %v", err)
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid endTime format: %v", err)
	}

	input := &cloudtrail.LookupEventsInput{
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	resp, err := a.Client.LookupEvents(ctx, input)
	if err != nil {
		fmt.Println("CloudTrail service unavailable. Using mock data instead.")
		return createMockEvents(), nil
	}
	return resp.Events, nil
}

func createMockEvents() []types.Event {
	events := make([]types.Event, 10)
	actionTypes := []string{"CreateUser", "DeleteUser", "UpdatePolicy", "AssumeRole", "ConsoleLogin",
		"PutObject", "GetObject", "CreateTable", "RunInstances", "StartInstance"}
	sources := []string{"iam.amazonaws.com", "s3.amazonaws.com", "ec2.amazonaws.com",
		"dynamodb.amazonaws.com", "cloudformation.amazonaws.com"}
	for i := 0; i < 10; i++ {
		eventId := fmt.Sprintf("mock-event-%d", i+1)
		username := fmt.Sprintf("user%d", (i%3)+1)
		eventName := actionTypes[i]
		eventSource := sources[i%len(sources)]
		eventTime := time.Now().Add(time.Duration(-i) * time.Hour)
		events[i] = types.Event{
			EventId:     &eventId,
			Username:    &username,
			EventName:   &eventName,
			EventSource: &eventSource,
			EventTime:   &eventTime,
		}
	}

	return events
}
