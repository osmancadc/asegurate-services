package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerInternalData(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		return ErrorMessage(err)
	}

	conn, err := ConnectDatabase()
	if err != nil {
		return ErrorMessage(err)
	}

	switch reqBody.Action {
	case `insertPerson`:
		return InsertPerson(conn, reqBody.PersonBody)
	case `updatePerson`:
		return UpdatePerson(conn, reqBody.PersonBody)
	case `getPersonByDocument`:
		return GetPersonByDocument(conn, reqBody.GetByDocumentBody)
	case `getNameByPhone`:
		return GetNameByPhone(conn, reqBody.GetByPhoneBody)
	case `insertUser`:
		return InsertUser(conn, reqBody.UserBody)
	case `updateUser`:
		return UpdateUser(conn, reqBody.UserBody)
	case `getUserByPhone`:
		return GetUserByPhone(conn, reqBody.GetByPhoneBody)
	case `checkUserByDocument`:
		return CheckUserByDocument(conn, reqBody.GetByDocumentBody)
	case `getAccountdata`:
		return GetAccountData(conn, reqBody.GetByDocumentBody)
	case `insertScore`:
		return InsertScore(conn, reqBody.ScoreBody)
	case `getScore`:
		return GetScoreByDocument(conn, reqBody.GetByDocumentBody)

	}

	response := SetResponseHeaders()
	response.StatusCode = http.StatusBadRequest
	response.Body = `{"message":"not a valid action"}`
	return response, nil
}

func main() {
	lambda.Start(HandlerInternalData)
}
