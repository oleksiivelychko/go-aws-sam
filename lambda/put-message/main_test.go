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
		resp, err := handler(event{})
		if err != nil && err.Error() != "got empty SQS name" {
			t.Errorf("got invalid HTTP newResponse: %v", resp)
		}
	})

	t.Run("Successful request", func(t *testing.T) {
		_, err := handler(event{Queue: "000000000000/MyQueue"})
		if err != nil {
			t.Error(err)
		}
	})
}
