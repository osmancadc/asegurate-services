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

	switch reqBody.Action {
	case `getPersonData`:
		return GetPersonData(reqBody.DataBody)
	case `getPersonName`:
		return GetPersonName(reqBody.NameBody)
	case `getPersonProccedings`:
		return GetProccedings(reqBody.ProccedingsBody)
	}

	response.StatusCode = http.StatusBadRequest
	response.Body = `{"message":"not a valid action"}`
	return response, nil

}

func main() {
	lambda.Start(HandlerExternalScoreData)
}
