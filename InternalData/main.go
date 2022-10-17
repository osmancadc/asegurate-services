package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Println(err.Error())
	}

	switch reqBody.Action {
	case `insertScore`:
		return InsertScore(conn, reqBody.ScoreBody)
	case `updatePerson`:
		return UpdatePerson(conn, reqBody.PersonBody)
	case `getScore`:
		return GetScoreByDocument(conn, reqBody.GetByDocumentBody)
	case `getUserByPhone`:
		return GetUserByPhone(conn, reqBody.GetByPhoneBody)
	case `getPersonByDocument`:
		return GetPersonByDocument(conn, reqBody.GetByDocumentBody)
	case `checkUserByDocument`:
		return CheckUserByDocument(conn, reqBody.GetByDocumentBody)
	case `insertUser`:
		return InsertUser(conn, reqBody.UserBody)
	case `insertPerson`:
		return InsertPerson(conn, reqBody.PersonBody)
	case `getAccountdata`:
		return GetAccountData(conn, reqBody.GetByDocumentBody)
	}

	response := SetResponseHeaders()
	response.StatusCode = http.StatusBadRequest
	response.Body = `{"message":"not a valid action"}`
	return response, nil
}

func main() {
	lambda.Start(HandlerInternalScoreData)
}
