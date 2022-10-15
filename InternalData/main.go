package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerInternalScoreData(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody

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
		return InsertInternalScore(conn, reqBody.InsertScoreBody)
	case `updateScore`:
		return UpdateInternalScore(conn, reqBody.UpdateScoreBody)
	case `getScore`:
		return GetInternalScoreSummary(conn, reqBody.GetScoreBody)
	case `getUserByPhone`:
		return GetUserByPhone(conn, reqBody.GetUserByPhoneBody)
	case `getPersonByDocument`:
		return GetPersonByDocument(conn, reqBody.GetByDocumentBody)
	case `getUserByDocument`:
		return GetUserByDocument(conn, reqBody.GetByDocumentBody)
	case `insertUser`:
		return InsertUser(conn, reqBody.InsertUserBody)
	case `insertPerson`:
		return InsertPerson(conn, reqBody.InsertPersonBody)
	}

	response := SetResponseHeaders()
	response.StatusCode = http.StatusBadRequest
	response.Body = `{"message":"not a valid action"}`
	return response, nil
}

func main() {
	lambda.Start(HandlerInternalScoreData)
}
