package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerAuthenticateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var data RequestBody

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	conn, err := ConnectDatabase()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}
	defer conn.Close()

	found, user, err := GetUserData(conn, data)
	if err != nil {
		response.Body = fmt.Sprintf(`{"message":"%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	if !found {
		response.Body = `{"message":"No user found","token":""}`
		response.StatusCode = http.StatusUnauthorized
		return response, nil
	}

	token, err := GenerateJWT(user)
	if err != nil {
		response.Body = fmt.Sprintf(`{"message":"%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	response.Body = fmt.Sprintf(`{"message":"User authenticated","token":"%s"}`, token)
	response.StatusCode = http.StatusOK
	return response, nil

}

func main() {
	lambda.Start(HandlerAuthenticateUser)
}
