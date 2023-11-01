package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type event struct {
	Queue string `json:"queue"`
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var e event
	err := json.Unmarshal([]byte(req.Body), &e)
	if err != nil {
		log.Printf("got invalid HTTP request %s: %s\n", req.Body, err.Error())
		return response(err.Error(), http.StatusInternalServerError, nil)
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Println("got empty AWS_REGION")
		return response("got empty AWS_REGION", http.StatusBadRequest, nil)
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		log.Println("got empty AWS_ACCESS_KEY_ID")
		return response("got empty AWS_ACCESS_KEY_ID", http.StatusBadRequest, nil)
	}

	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		log.Println("got empty AWS_SECRET_ACCESS_KEY")
		return response("got empty AWS_SECRET_ACCESS_KEY", http.StatusBadRequest, nil)
	}

	endpoint := os.Getenv("AWS_ENDPOINT")
	if endpoint == "" {
		log.Println("got empty AWS_ENDPOINT")
		return response("got empty AWS_ENDPOINT", http.StatusBadRequest, nil)
	}

	config := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	}

	if endpoint != "" {
		config.Endpoint = aws.String(endpoint)
	}

	awsSession, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}

	if e.Queue == "" {
		log.Println("got empty queue name")
		return response("", http.StatusBadRequest, errors.New("got empty queue name"))
	}

	sqsSession := sqs.New(session.Must(awsSession, nil))
	queueURL := fmt.Sprintf("%s/%s", endpoint, e.Queue)

	err = sendToQueue(sqsSession, queueURL)
	if err != nil {
		log.Printf("unable to put message: %s\n", err)
		return response("unable to put message", http.StatusInternalServerError, err)
	}

	log.Printf("successfully put message into %s\n", queueURL)
	return response(fmt.Sprintf("successfully put message into %s\n", queueURL), http.StatusOK, nil)
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

func response(body string, code int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: code,
	}, err
}
