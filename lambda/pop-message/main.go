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

func main() { lambda.Start(handler) }

func printMessage(m events.SQSMessage) {
	if attr, ok := m.MessageAttributes["MyAttr"]; ok {
		fmt.Printf(
			"\nMessage %q of event source %q contains body %q, 'MyAttr' %q\n",
			m.MessageId,
			m.EventSource,
			m.Body,
			*attr.StringValue,
		)
		return
	}

	fmt.Printf(
		"\nMessage %q of event source %q contains body %q\n",
		m.MessageId,
		m.EventSource,
		m.Body,
	)
}
