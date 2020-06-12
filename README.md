# AWS Lambda With Golang

This project was created to test AWS Lambda with Golang sending message to SQS and triggered by an API Gateway

## Getting Started

* Create a lambda function and choose the language Golang
* To deploy:
  - GOOS=linux go build -o main
  - zip main.zip main
  - Upload the zip file to S3 and link it to the lambda function
* Create an API Gateway, link the lambda function to the method http action and check "lambda proxy" checkbox
* Create a SQS
* Create the environment variable in the lambda with the "REGION" and "SQS_URL"

### Prerequisites

```
Golang 1.14
AWS
```
