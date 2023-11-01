package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type event struct {
	Queue string `json:"queue"`
}

type response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

var (
	awsRegion          string
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsEndpoint        string
)

func handler(e event) (response, error) {
	if e.Queue == "" {
		log.Println("got empty queue name")
		return newResponse("", http.StatusBadRequest, errors.New("got empty queue name"))
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
	queueURL := fmt.Sprintf("%s/%s", awsEndpoint, e.Queue)

	err = sendToQueue(sqsSession, queueURL)
	if err != nil {
		log.Printf("unable to put message: %s\n", err)
		return newResponse("unable to put message", http.StatusInternalServerError, err)
	}

	log.Printf("successfully put message into %s\n", queueURL)
	return newResponse(fmt.Sprintf("successfully put message into %s", queueURL), http.StatusOK, nil)
}

func main() {
	lambda.Start(handler)
}

func sendToQueue(sqsSession *sqs.SQS, url string) error {
	_, err := sqsSession.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"MyAttr": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("Time now is %s", time.Now().Format(time.DateTime))),
			},
		},
		MessageBody: aws.String("Got new event!"),
		QueueUrl:    aws.String(url),
	})

	return err
}

func newResponse(body string, code int, err error) (response, error) {
	return response{
		Body:       body,
		StatusCode: code,
	}, err
}
