package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var awsRegion = os.Getenv("REGION")
var quereURL = os.Getenv("SQS_URL")
var svc = sqs.New(session.New(), aws.NewConfig().WithRegion(awsRegion))
var jokeEndpoint = "https://us-central1-kivson.cloudfunctions.net/charada-aleatoria"
var client = resty.New()

func main() {
	lambda.Start(router)
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Path == "/joke" {
		if req.HTTPMethod == "GET" {
			resp, err := client.R().SetHeader("Accept", "application/json").Get(jokeEndpoint)
			if err != nil {
				return errorAPIGatewayResponse(http.StatusNotFound, err.Error())
			}

			joke := &joke{}
			json.Unmarshal(resp.Body(), joke)

			insertDataToSqs(string(resp.Body()))
			// time.Sleep(3 * time.Second)

			// response, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			// 	QueueUrl:              aws.String(quereURL),
			// 	MessageAttributeNames: []*string{aws.String("All")},
			// 	MaxNumberOfMessages:   aws.Int64(10),
			// })
			// fmt.Println(response.Messages)

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
	attribute := &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String("teste"),
	}
	mapAttributes := make(map[string]*sqs.MessageAttributeValue)
	mapAttributes["key"] = attribute

	sqsMessage := &sqs.SendMessageInput{
		MessageBody:       aws.String(message),
		QueueUrl:          aws.String(quereURL),
		MessageAttributes: mapAttributes,
	}

	svc.SendMessage(sqsMessage)
}

type joke struct {
	ID       int    `json:id`
	Pergunta string `json:pergunta`
	Resposta string `json:resposta`
}
