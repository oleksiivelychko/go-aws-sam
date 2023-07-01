include ../.env

sam-build:
	sam build

sam-invoke: sam-build
	sam local invoke -e events/event.json

sam-start:
	sam local start-api

zip:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto \
		go build -C lambda \
		-ldflags "-X main.awsRegion=$(AWS_REGION) -X main.awsAccessKeyID=$(AWS_ACCESS_KEY_ID) -X main.awsSecretAccessKey=$(AWS_SECRET_ACCESS_KEY) -X main.awsEndpoint=$(AWS_ENDPOINT)" \
		-o ../handler-bin
	zip lambda.zip handler-bin

list:
	aws lambda list-functions --endpoint-url=http://localhost:4566 --profile localstack

create:
	aws lambda create-function --function-name put-message \
        --runtime go1.x \
        --zip-file fileb://lambda.zip \
        --handler handler-bin \
        --memory-size 128 \
        --role arn:aws:iam::000000000000:role/lambda-role \
    	--endpoint-url=http://localhost:4566 --profile localstack

update: zip
	aws lambda update-function-code --function-name put-message \
        --zip-file fileb://lambda.zip \
    	--endpoint-url=http://localhost:4566 --profile localstack

create-config:
	aws lambda create-function-url-config --function-name put-message \
        --auth-type NONE \
        --endpoint-url=http://localhost:4566 --profile localstack

status:
	aws lambda get-function --function-name put-message --endpoint-url=http://localhost:4566 --profile localstack

delete:
	aws lambda delete-function --function-name put-message --endpoint-url=http://localhost:4566 --profile localstack

delete-config:
	aws lambda delete-function-url-config --function-name put-message --endpoint-url=http://localhost:4566 --profile localstack

invoke:
	aws lambda invoke --function-name put-message \
		--cli-binary-format raw-in-base64-out \
		--payload '{"queue": "000000000000/my-queue"}'
		--endpoint-url=http://localhost:4566 --profile localstack outfile.json

post:
	curl -X POST \
    	'FUNCTION_URL' \
    	-H 'Content-Type: application/json' \
    	-d '{"queue": "000000000000/my-queue"}'
