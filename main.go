package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"io"
	"net/http"
)

type InputEvent struct {
	Bucket   string `json:"bucket"`
	Url      string `json:"url"`
	SaveAs   string `json:"saveAs"`
	QueueUrl string `json:"queueUrl"`
}

type Request struct {
	Bucket   string `json:"bucket"`
	Url      string `json:"url"`
	SaveAs   string `json:"saveAs"`
	QueueUrl string `json:"queueUrl"`
	Body     string `json:"body"`
}

type Response struct {
	Content    string   `json:"body"`
	Headers    []string `json:"headers"`
	StatusCode int      `json:"statusCode"`
}

var s3session *s3.S3
var sqsSession *sqs.SQS
var awsRegion string
var awsAccessKeyId string
var awsSecretAccessKeyId string

func init() {
	if awsRegion == "" || awsAccessKeyId == "" || awsSecretAccessKeyId == "" {
		panic("check environment variables: some were skipped during the build phase")
	}

	config := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKeyId, ""),
	}

	session, err := awsSession.NewSession(config)
	if err != nil {
		panic(err)
	}

	s3session = s3.New(awsSession.Must(session, nil))
	sqsSession = sqs.New(awsSession.Must(session, nil))
}

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(inputEvent Request) (Response, error) {
	bucket := inputEvent.Bucket
	url := inputEvent.Url
	saveAs := inputEvent.SaveAs
	queueUrl := inputEvent.QueueUrl

	// data came through API Gateway POST request
	if inputEvent.Body != "" {
		input := &InputEvent{}
		_ = json.Unmarshal([]byte(inputEvent.Body), &input)

		bucket = input.Bucket
		url = input.Url
		saveAs = input.SaveAs
		queueUrl = input.QueueUrl
	}

	if bucket == "" || url == "" || saveAs == "" {
		return getResponse(fmt.Sprintf("bucket: %s, url: %s, saveAs: %s", bucket, url, saveAs), http.StatusBadRequest, errors.New("required attributes are missing"))
	}

	fileBytes, err := getFile(url)
	if err != nil {
		return getResponse(fmt.Sprintf("unable to get file from '%s'", url), http.StatusBadRequest, err)
	}

	_, err = s3session.PutObject(&s3.PutObjectInput{
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
		Body:   bytes.NewReader(fileBytes),
		Bucket: aws.String(bucket),
		Key:    aws.String(saveAs),
	})

	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			err = errors.New(fmt.Sprintf(
				"AWS error: %s\noriginal: %s\ncode: %s\nmessage: %s",
				awsError.Error(), awsError.OrigErr(), awsError.Code(), awsError.Message()))
		}
		return getResponse("caught AWS error", http.StatusBadRequest, err)
	}

	s3Link := fmt.Sprintf("https://%s.%s.amazonaws.com/%s", bucket, awsRegion, saveAs)

	if queueUrl != "" {
		err = sendToQueue(url, s3Link)
		if err != nil {
			return getResponse("unable to put message into queue", http.StatusInternalServerError, err)
		}
	}

	return getResponse(s3Link, http.StatusOK, nil)
}

func sendToQueue(link, s3Link string) error {
	_, err := sqsSession.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Link": {
				DataType:    aws.String("String"),
				StringValue: aws.String(link),
			},
			"S3Link": {
				DataType:    aws.String("String"),
				StringValue: aws.String(s3Link),
			},
		},
	})

	return err
}

func getFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	fileBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func getResponse(content string, status int, err error) (Response, error) {
	return Response{
		Content:    content,
		Headers:    []string{"Accept: application/json"},
		StatusCode: status,
	}, err
}
