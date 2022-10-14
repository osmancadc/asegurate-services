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
		return ErrorMessage(err)
	}

	conn, err := ConnectDatabase()
	if err != nil {
		return ErrorMessage(err)
	}
	defer conn.Close()

	switch reqBody.Action {
	case `insertScore`:
		return UploadInternalScore(conn, reqBody.InsertScoreBody)
	case `updateScore`:
		return UpdateInternalScore(conn, reqBody.UpdateScoreBody)
	case `getScore`:
		return GetInternalScoreSummary(conn, reqBody.GetScoreBody)
	case `getByPhone`:
		return GetUserByPhone(conn, reqBody.GetByPhoneBody)
	}

	response.StatusCode = http.StatusBadRequest
	response.Body = `{"message":"not a valid action"}`
	return response, nil
}

func main() {
	lambda.Start(HandlerInternalScoreData)
}
