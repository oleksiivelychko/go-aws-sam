package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) {
	for _, message := range sqsEvent.Records {
		printMessage(message)
	}
}

func main() {
	lambda.Start(handler)
}

func printMessage(message events.SQSMessage) {
	fmt.Printf(
		"\nMessage ID %s for event source %s contains body %s with attribute `MyAttr`=%s\n",
		message.MessageId,
		message.EventSource,
		message.Body,
		*message.MessageAttributes["MyAttr"].StringValue,
	)
}
