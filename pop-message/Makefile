build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=auto go build -C lambda -o ../handler-bin
	zip lambda.zip handler-bin

create: build
	aws lambda create-function --function-name pop-message \
        --runtime go1.x \
        --zip-file fileb://lambda.zip \
        --handler handler-bin \
        --memory-size 128 \
        --role arn:aws:iam::000000000000:role/lambda-role \
    	--endpoint-url=http://localhost:4566 --profile localstack

mapping:
	aws lambda create-event-source-mapping --function-name pop-message \
        --event-source-arn arn:aws:sqs:us-east-1:000000000000:my-queue \
        --endpoint-url=http://localhost:4566 --profile localstack

put-rule:
	aws events put-rule --name api-rule \
        --schedule-expression 'rate(1 minute)' \
        --endpoint-url=http://localhost:4566 --profile localstack

put-targets:
	aws events put-targets --rule api-rule \
        --targets '[{"Id":"1","Arn":"arn:aws:sqs:us-east-1:000000000000:my-queue","Input":"{\"MyAttr\":\"test\"}"}]' \
        --endpoint-url=http://localhost:4566 --profile localstack

send:
	aws sqs send-message --queue-url http://localhost:4566/000000000000/my-queue \
		--message-body "Hello, World!" \
		--endpoint-url=http://localhost:4566 --profile localstack

logs:
	aws logs tail /aws/lambda/pop-message --follow --endpoint-url=http://localhost:4566 --profile localstack


