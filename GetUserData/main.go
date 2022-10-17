package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerGetUserData(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	document := req.PathParameters["document"]

	client := GetClient()

	user, err := GetUserData(document, client)
	if err != nil {
		return ErrorMessage(err)
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return ErrorMessage(err)
	}

	response := SetResponseHeaders()
	response.Body = string(userJson)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HandlerGetUserData)
}
