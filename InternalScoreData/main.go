package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerInternalScoreData(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
		return response, nil
	}

	conn, err := ConnectDatabase()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}
	defer conn.Close()

	switch reqBody.Action {
	case `insert`:
		return InsertInternalScore(conn, reqBody.InsertBody)
	case `update`:
		return UpdateInternalScore(conn, reqBody.UpdateBody)
	case `get`:
		return GetInternalScoreSummary(conn, reqBody.GetBody)
	}

	response.StatusCode = http.StatusBadRequest
	response.Body = `{"message":"not a valid action"}`
	return response, nil
}

func main() {
	lambda.Start(HandlerInternalScoreData)
}
