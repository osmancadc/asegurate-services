package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerGetPersonName(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var reqBody RequestBody

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		response := SetResponseHeaders()
		response.StatusCode = http.StatusBadRequest
		return response, nil
	}

	client := GetClient()

	document := reqBody.Document

	email, err := GetEmailByDocument(document, client)
	if err != nil {
		return ErrorMessage(errors.New(`el usuario no se pudo encontrar`))
	}

	return SuccessMessage(email)

}

func main() {
	lambda.Start(HandlerGetPersonName)
}
