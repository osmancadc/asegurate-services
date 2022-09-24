package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HanderGetUserData(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	fmt.Println("====================================")
	fmt.Println(req.PathParameters["user_id"])
	fmt.Println("====================================")
	fmt.Printf("-> %v", req.PathParameters)
	fmt.Println("====================================")

	response.Body = req.QueryStringParameters["user_id"]
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderGetUserData)
}
