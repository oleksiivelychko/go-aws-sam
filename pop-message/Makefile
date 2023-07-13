sam-build:
	sam build

sam-invoke: sam-build
	sam local invoke

zip:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto go build -C lambda -o ../handler-bin
	zip pop-message.zip handler-bin

logs:
	aws logs tail /aws/lambda/pop-message --follow --endpoint-url=http://localhost:4566 --profile localstack
