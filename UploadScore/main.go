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

	if reqBody.Type == `CC` {
		err = UploadScoreDocument(conn, reqBody.Author, reqBody.Score, reqBody.Value, reqBody.Comments)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}
	} else {
		err = UploadScorePhone(conn, reqBody.Author, reqBody.Score, reqBody.Value, reqBody.Comments)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}
	}

	response.Body = `{"message":"Score uploaded successfully"}`
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderUploadScore)
}
