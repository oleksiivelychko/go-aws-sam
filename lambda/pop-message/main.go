package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, e events.SQSEvent) {
	for _, m := range e.Records {
		printMessage(m)
	}
}

func main() {
	lambda.Start(handler)
}

func printMessage(m events.SQSMessage) {
	attr, ok := m.MessageAttributes["MyAttr"]
	if ok {
		fmt.Printf(
			"\nMessage ID %s for event source %s contains body `%s` with attribute 'MyAttr'=`%s`\n",
			m.MessageId,
			m.EventSource,
			m.Body,
			*attr.StringValue,
		)
		return
	}

	fmt.Printf(
		"\nMessage ID %s for event source %s contains body `%s`\n",
		m.MessageId,
		m.EventSource,
		m.Body,
	)
}
