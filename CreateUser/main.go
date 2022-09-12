package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HanderUploadScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	conn := ConnectDatabase()
	defer conn.Close()

	err = InsertPerson(conn, reqBody.Document, reqBody.Name, reqBody.Lastname)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	id, err := InsertUser(conn, reqBody.Username, reqBody.Email, reqBody.Phone, reqBody.Password, reqBody.Document, reqBody.Role)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if reqBody.Role == `seller` {
		err = InsertSeller(conn, id)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}
	}

	response.Body = fmt.Sprintf(`{ "message": "user created successfully (AWS)"}`)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderUploadScore)
}
