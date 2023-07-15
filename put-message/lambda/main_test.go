package main

import (
	"testing"
)

func TestHandler(t *testing.T) {
	awsRegion = "us-east-1"
	awsAccessKeyID = "local"
	awsSecretAccessKey = "local"
	awsEndpoint = "http://localhost:4566"

	t.Run("Failed request", func(t *testing.T) {
		resp, err := handler(Event{})
		if err != nil && err.Error() != "got empty SQS name" {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", resp)
		}
	})

	t.Run("Successful request", func(t *testing.T) {
		_, err := handler(Event{Queue: "000000000000/my-queue"})
		if err != nil {
			t.Fatalf("Everything should be ok, but %s", err)
		}
	})
}
