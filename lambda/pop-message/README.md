### Run **PopMessage** function.

- build function:
```
sam build PopMessage
```
- invoke function:
```
sam local generate-event sqs receive-message --body 'Hello, World!' | sam local invoke -e - PopMessage
```
---

- debug function (`delve` package must be installed):
```
sam local invoke PopMessage -e events/event.json -d 2345 --debugger-path=delve --debug-args="-delveAPI=2" --debug
```
---

- show logs:
```
aws logs tail /aws/lambda/pop-message --follow --endpoint-url=http://localhost:4566 --profile localstack
```
---

- create IDE configuration to run/debug:
![SAM IDE configuration](sam_ide_run_configuration.png)
---
