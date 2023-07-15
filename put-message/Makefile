include ../.env

invoke:
	sam build
	sam local invoke --payload events/event.json

zip:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto \
		go build -C lambda \
		-ldflags "-X main.awsRegion=$(AWS_REGION) -X main.awsAccessKeyID=$(AWS_ACCESS_KEY_ID) -X main.awsSecretAccessKey=$(AWS_SECRET_ACCESS_KEY) -X main.awsEndpoint=$(AWS_ENDPOINT)" \
		-o ../handler-bin
	zip put-message.zip handler-bin
