package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HanderGetScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	score, isStored, err := GetStoredScore(conn, reqBody.Document)

	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if !isStored {
		fmt.Println("No internal data was found")

		score, err := CalculateScore(reqBody.Document, reqBody.Type)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		err = SaveNewPerson(conn, score, reqBody.Document)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		response.Body = GetResponseBody(score, reqBody.Document)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	fmt.Println("Internal data found")

	elapsed, err := DaysSinceLastUpdate(score.Updated)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if elapsed > 7 {
		fmt.Println("Updated over a week ago")
		score, _ := CalculateScore(reqBody.Document, reqBody.Type)

		response.Body = GetResponseBody(score, reqBody.Document)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	response.Body = GetResponseBody(score, reqBody.Document)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderGetScore)
}
