package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"net/http"
	"os"
	"time"
)

type event struct {
	Queue string `json:"queue"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var e event
	err := json.Unmarshal([]byte(request.Body), &e)
	if err != nil {
		return response(err.Error(), http.StatusInternalServerError, nil)
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		return response("got empty AWS_REGION", http.StatusBadRequest, nil)
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKeyID == "" {
		return response("got empty AWS_ACCESS_KEY_ID", http.StatusBadRequest, nil)
	}

	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretAccessKey == "" {
		return response("got empty AWS_SECRET_ACCESS_KEY", http.StatusBadRequest, nil)
	}

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint == "" {
		return response("got empty AWS_ENDPOINT", http.StatusBadRequest, nil)
	}

	awsCfg := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	}

	if awsEndpoint != "" {
		awsCfg.Endpoint = aws.String(awsEndpoint)
	}

	awsSession, err := session.NewSession(awsCfg)
	if err != nil {
		panic(err)
	}

	sqsSession := sqs.New(session.Must(awsSession, nil))

	if e.Queue == "" {
		return response("", http.StatusBadRequest, errors.New("got empty SQS name"))
	}

	queueURL := fmt.Sprintf("%s/%s", awsEndpoint, e.Queue)

	err = sendToQueue(sqsSession, queueURL)
	if err != nil {
		return response("unable to put message", http.StatusInternalServerError, err)
	}

	return response(fmt.Sprintf("successfully put message by QueueUrl: %s\n", queueURL), http.StatusOK, nil)
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

func response(body string, statusCode int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: statusCode,
	}, err
}
