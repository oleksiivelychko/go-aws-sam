package main

import (
	"github.com/aws/aws-lambda-go/events"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	_ = os.Setenv("AWS_REGION", "us-east-1")
	_ = os.Setenv("AWS_ACCESS_KEY_ID", "local")
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "local")
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "http://localhost:4566")

	t.Run("Failed request", func(t *testing.T) {
		resp, err := handler(events.APIGatewayProxyRequest{})
		if err != nil && err.Error() != "unexpected end of JSON input" {
			t.Errorf("Error failed to trigger with an invalid HTTP response: %v", resp)
		}
	})

	t.Run("Successful request", func(t *testing.T) {
		_, err := handler(events.APIGatewayProxyRequest{
			Body: "{\"queue\":\"000000000000/MyQueue\"}",
		})

		if err != nil {
			t.Error(err)
		}
	})
}
