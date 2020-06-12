package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var awsRegion = os.Getenv("REGION")
var quereURL = os.Getenv("SQS_URL")
var svc = sqs.New(session.New(), aws.NewConfig().WithRegion(awsRegion))

func main() {
	fmt.Printf("Generated service")
	lambda.Start(router)
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Path == "/users" {
		if req.HTTPMethod == "GET" {
			message := fmt.Sprintf("This is only a register of a GET request (%s)", time.Now())
			insertDataToSqs(message)
			return successAPIGatewayResponse()
		}

		if req.HTTPMethod == "POST" {
			var user user
			err := json.Unmarshal([]byte(req.Body), &user)

			if err != nil {
				return errorAPIGatewayResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			}

			message := fmt.Sprintf("Sending post message to the lambda with user %s", user.Name)
			insertDataToSqs(message)

			return successAPIGatewayResponse()
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       http.StatusText(http.StatusMethodNotAllowed),
	}, nil

}

func successAPIGatewayResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "{ \"status,\": \"200\" }",
	}, nil
}

func errorAPIGatewayResponse(status int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
	}, nil
}

func insertDataToSqs(message string) {
	fmt.Println("Inserting data to sqs")
	sqsMessage := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(quereURL),
	}

	svc.SendMessage(sqsMessage)
}

type user struct {
	Name  string `json:Name`
	email string `json:email`
}
