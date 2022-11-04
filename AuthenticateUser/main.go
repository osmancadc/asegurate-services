package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerAuthenticateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		return ErrorMessage(err, http.StatusInternalServerError)
	}

	client := GetClient()

	found, valid, err := ValidateUser(reqBody, client)
	if err != nil {
		return ErrorMessage(err, http.StatusInternalServerError)
	}

	if !found {
		return ErrorMessage(errors.New(`not user found`), http.StatusNotFound)
	}

	if !valid {
		return ErrorMessage(errors.New(`incorrect password`), http.StatusUnauthorized)
	}

	token, err := GenerateJWT(reqBody.Document)
	if err != nil {
		return ErrorMessage(err, http.StatusInternalServerError)
	}

	return SuccessMessage(token)
}

func main() {
	lambda.Start(HandlerAuthenticateUser)
}
