### Run **PutMessage** function.

⚠️ Queue **MyQueue** must be created before.

- build function:
```
sam build PutMessage
```
- invoke function:
```
sam local invoke PutMessage -e events/event.json --skip-pull-image
```
    
