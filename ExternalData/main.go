package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerExternalScoreData(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody
	response := SetResponseHeaders()

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
	case `getPersonData`:
		return GetPersonData(reqBody.DataBody)
	case `getPersonName`:
		return GetPersonName(reqBody.NameBody)
	}

	response.StatusCode = http.StatusBadRequest
	response.Body = `{"message":"not a valid action"}`
	return response, nil

}

func main() {
	lambda.Start(HandlerExternalScoreData)
}
