package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerCreateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody
	response := SetResponseHeaders()

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, nil
	}

	client := GetClient()

	name, err := InsertPerson(reqBody.Document, reqBody.ExpeditionDate, client)
	if err != nil {
		return ErrorMessage(err)
	}

	user := UserBody{
		Document: reqBody.Document,
		Email:    reqBody.Email,
		Phone:    reqBody.Phone,
		Password: reqBody.Password,
		Role:     reqBody.Role,
	}

	err = InsertUser(user, client)
	if err != nil {
		return ErrorMessage(err)
	}

	response.Body = fmt.Sprintf(`{ "message": "user created successfully","name":"%s"}`, name)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HandlerCreateUser)
}
