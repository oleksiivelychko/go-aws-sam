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
	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_REGION")
		if awsRegion == "" {
			log.Print("empty AWS_REGION")
			return newResponse(fmt.Sprintf("empty AWS_REGION"), http.StatusBadRequest, errors.New("empty AWS_REGION"))
		}
	}

	if awsAccessKeyID == "" {
		awsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		if awsAccessKeyID == "" {
			log.Print("empty AWS_ACCESS_KEY_ID")
			return newResponse(fmt.Sprintf("empty AWS_ACCESS_KEY_ID"), http.StatusBadRequest, errors.New("empty AWS_ACCESS_KEY_ID"))
		}
	}

	if awsSecretAccessKey == "" {
		awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		if awsSecretAccessKey == "" {
			log.Print("empty AWS_SECRET_ACCESS_KEY")
			return newResponse(fmt.Sprintf("empty AWS_SECRET_ACCESS_KEY"), http.StatusBadRequest, errors.New("empty AWS_SECRET_ACCESS_KEY"))
		}
	}

	if awsEndpoint == "" {
		awsEndpoint = os.Getenv("AWS_ENDPOINT")
		if awsEndpoint == "" {
			log.Print("empty AWS_ENDPOINT")
			return newResponse(fmt.Sprintf("empty AWS_ENDPOINT"), http.StatusBadRequest, errors.New("empty AWS_ENDPOINT"))
		}
	}

	if e.Queue == "" {
		log.Println("empty queue parameter")
		return newResponse(fmt.Sprintf("empty queue parameter"), http.StatusBadRequest, errors.New("empty queue parameter"))
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
		return newResponse(err.Error(), http.StatusInternalServerError, err)
	}

	queueURL := fmt.Sprintf("%s/%s", awsEndpoint, e.Queue)

	err = sendToQueue(sqs.New(session.Must(awsSession, nil)), queueURL)
	if err != nil {
		log.Print(err)
		return newResponse(fmt.Sprintf("%s", err), http.StatusInternalServerError, err)
	}

	log.Printf("message was put into %s", queueURL)
	return newResponse(fmt.Sprintf("message was put into %s", queueURL), http.StatusOK, nil)
}

func main() { lambda.Start(handler) }

func sendToQueue(sqsSession *sqs.SQS, url string) error {
	_, err := sqsSession.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"MyAttr": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("Time is %s", time.Now().Format(time.DateTime))),
			},
		},
		MessageBody: aws.String("New event!"),
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
