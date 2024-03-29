package main

import (
	"testing"
)

func TestHandler(t *testing.T) {
	awsRegion = "us-east-1"
	awsAccessKeyID = "local"
	awsSecretAccessKey = "local"
	awsEndpoint = "http://localhost:4566"

	t.Run("Failed HTTP request", func(t *testing.T) {
		if resp, err := handler(&event{}); err != nil && err.Error() != "empty SQS" {
			t.Errorf("bad response: %v", resp)
		}
	})

	t.Run("Successful HTTP request", func(t *testing.T) {
		if _, err := handler(&event{Queue: "000000000000/MyQueue"}); err != nil {
			t.Error(err)
		}
	})
}
