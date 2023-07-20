package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Queue string `json:"queue"`
}

type Response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

var (
	awsRegion          string
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsEndpoint        string
)

func handler(event Event) (Response, error) {
	if event.Queue == "" {
		return response("", http.StatusBadRequest, errors.New("got empty SQS name"))
	}

	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_REGION")
		if awsRegion == "" {
			panic("got empty AWS_REGION")
		}
	}
	if awsAccessKeyID == "" {
		awsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		if awsAccessKeyID == "" {
			panic("got empty AWS_ACCESS_KEY_ID")
		}
	}
	if awsSecretAccessKey == "" {
		awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		if awsSecretAccessKey == "" {
			panic("got empty AWS_SECRET_ACCESS_KEY")
		}
	}
	if awsEndpoint == "" {
		awsEndpoint = os.Getenv("AWS_ENDPOINT")
		if awsEndpoint == "" {
			panic("got empty AWS_ENDPOINT")
		}
	}

	config := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	}

	if awsEndpoint != "" {
		config.Endpoint = aws.String(awsEndpoint)
	}

	awsSession, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}

	sqsSession := sqs.New(session.Must(awsSession, nil))

	queueURL := fmt.Sprintf("%s/%s", awsEndpoint, event.Queue)

	err = sendToQueue(sqsSession, queueURL)
	if err != nil {
		return response("unable to put message", http.StatusInternalServerError, err)
	}

	return response(fmt.Sprintf("successfully put message by QueueUrl: %s", queueURL), http.StatusOK, nil)
}

func main() {
	lambda.Start(handler)
}

func sendToQueue(sqsSession *sqs.SQS, queueURL string) error {
	_, err := sqsSession.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"MyAttr": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("Time now is %s", time.Now().Format(time.DateTime))),
			},
		},
		MessageBody: aws.String("Got new event!"),
		QueueUrl:    aws.String(queueURL),
	})

	return err
}

func response(content string, status int, err error) (Response, error) {
	return Response{
		Body:       content,
		StatusCode: status,
	}, err
}
