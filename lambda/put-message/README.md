### Run **PutMessage** function.

- build function:
```
sam build PutMessage
```
- invoke function:
```
sam local invoke PutMessage -e events/event.json
```
---

- running `PutMessage` function as (development) server:
```
sam local invoke PutMessage -e events/event.json
```
- [invoke function](https://github.com/oleksiivelychko/go-aws-cli/blob/main/md/lambda.md#local-usage-of-aws-lambda-via-cobra-andor-aws-cli) via AWS CLI
---

- build binary (make environment variables available from Makefile by `include ../.env`):
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto \
	go build -C lambda/put-message \
	-ldflags "-X main.awsRegion=$(AWS_REGION) -X main.awsAccessKeyID=$(AWS_ACCESS_KEY_ID) -X main.awsSecretAccessKey=$(AWS_SECRET_ACCESS_KEY) -X main.awsEndpoint=$(AWS_ENDPOINT)" \
	-o handler-bin
```
- create ZIP archive:
```
zip put-message.zip handler-bin
```
