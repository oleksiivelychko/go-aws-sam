package main

import (
	"encoding/json"
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
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Print("empty AWS_REGION")
		return response(fmt.Sprintf("empty AWS_REGION\n"), http.StatusBadRequest)
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		log.Print("empty AWS_ACCESS_KEY_ID")
		return response(fmt.Sprintf("empty AWS_ACCESS_KEY_ID\n"), http.StatusBadRequest)
	}

	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		log.Print("empty AWS_SECRET_ACCESS_KEY")
		return response(fmt.Sprintf("empty AWS_SECRET_ACCESS_KEY\n"), http.StatusBadRequest)
	}

	endpoint := os.Getenv("AWS_ENDPOINT")
	if endpoint == "" {
		log.Print("empty AWS_ENDPOINT")
		return response(fmt.Sprintf("empty AWS_ENDPOINT\n"), http.StatusBadRequest)
	}

	var e event

	if err := json.Unmarshal([]byte(req.Body), &e); err != nil {
		log.Printf("%s: %s", err, req.Body)
		return response(fmt.Sprintf("%s: %s\n", err, req.Body), http.StatusBadRequest)
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
		log.Print("empty queue parameter")
		return response(fmt.Sprintf("empty queue parameter\n"), http.StatusBadRequest)
	}

	queueURL := fmt.Sprintf("%s/%s", endpoint, e.Queue)

	err = sendToQueue(sqs.New(session.Must(awsSession, nil)), queueURL)
	if err != nil {
		log.Print(err.Error())
		return response(fmt.Sprintf("%s\n", err), http.StatusInternalServerError)
	}

	log.Printf("message was put into %s", queueURL)
	return response(fmt.Sprintf("message was put into %s\n", queueURL), http.StatusOK)
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

func response(body string, code int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: code,
	}, nil
}
