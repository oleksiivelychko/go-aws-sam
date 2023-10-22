# go-aws-sam

### Running AWS Serverless Application Model (SAM).

⚠️ SQS queue must be created before:
```
aws sqs create-queue --queue-name MyQueue --endpoint-url http://localhost:4566 --profile localstack
```

📌 [Run **PutMessageApi** function](lambda/put-message-api/README.md)

📌 [Run **PutMessage** function](lambda/put-message/README.md)

📌 [Run **PopMessage** function](lambda/pop-message/README.md)

📎 Create a new SAM configuration:
```
sam init --runtime go1.x
```

🐞 Debugging SAM on Apple M1 has bug https://github.com/aws/aws-toolkit-jetbrains/issues/3061
