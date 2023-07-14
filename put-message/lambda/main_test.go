package main

import (
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestHandler(t *testing.T) {
	awsRegion = "us-east-1"
	awsAccessKeyID = "local"
	awsSecretAccessKey = "local"
	awsEndpoint = "http://localhost:4566"

	t.Run("Failed request", func(t *testing.T) {
		resp, err := handler(events.APIGatewayProxyRequest{})
		if err != nil && err.Error() != "unexpected end of JSON input" {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", resp)
		}
	})

	t.Run("Successful request", func(t *testing.T) {
		_, err := handler(events.APIGatewayProxyRequest{
			Body: "{\"queue\":\"000000000000/my-queue\"}",
		})

		if err != nil {
			t.Fatalf("Everything should be ok, but %s", err)
		}
	})
}
