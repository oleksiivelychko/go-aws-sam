### Run **PutMessageApi** function.

⚠️ SQS must be created before.

- build function:
```
sam build PutMessageApi
```
- start local API:
```
sam local start-api
```
- send request to invoke function:
```
curl -X POST -d '{"queue":"000000000000/MyQueue"}' http://127.0.0.1:3000/api/put-message
```
