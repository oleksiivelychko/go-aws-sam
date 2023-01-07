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
	"io"
	"net/http"
)

type InputEvent struct {
	Bucket string `json:"bucket"`
	Url    string `json:"url"`
	SaveAs string `json:"saveAs"`
}

type Request struct {
	Bucket string `json:"bucket"`
	Url    string `json:"url"`
	SaveAs string `json:"saveAs"`
	Body   string `json:"body"`
}

type Response struct {
	Content    string   `json:"body"`
	Headers    []string `json:"headers"`
	StatusCode int      `json:"statusCode"`
}

var s3session *s3.S3
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
}

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(inputEvent Request) (Response, error) {
	bucket := inputEvent.Bucket
	url := inputEvent.Url
	saveAs := inputEvent.SaveAs

	// data came through API Gateway POST request
	if inputEvent.Body != "" {
		input := &InputEvent{}
		_ = json.Unmarshal([]byte(inputEvent.Body), &input)

		bucket = input.Bucket
		url = input.Url
		saveAs = input.SaveAs
	}

	if bucket == "" || url == "" || saveAs == "" {
		return getResponse(fmt.Sprintf("bucket: %s, url: %s, saveAs: %s", bucket, url, saveAs), http.StatusBadRequest, errors.New("required attributes are missing"))
	}

	fileBytes, err := getFile(url)
	if err != nil {
		return getResponse(fmt.Sprintf("unable to get file from '%s'", url), http.StatusBadRequest, err)
	}

	output, err := s3session.PutObject(&s3.PutObjectInput{
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

	return getResponse(output.String(), http.StatusOK, nil)
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
