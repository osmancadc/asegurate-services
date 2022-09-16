package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HanderGetScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req RequestBody

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	err := json.Unmarshal([]byte(req.Body), &req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	conn := ConnectDatabase()
	defer conn.Close()

	score, isStored, err := GetStoredScore(conn, req.Document)

	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if !isStored {
		fmt.Println("No se encontraron datos")

		score, _ := CalculateScore(req.Document, req.Type)

		response.Body = GetResponseBody(score, req.Document)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	fmt.Println("Se encontraron datos internos")

	elapsed, err := DaysSinceLastUpdate(score.Updated)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if elapsed > 7 {
		fmt.Println("Updated a week ago")
		score, _ := CalculateScore(req.Document, req.Type)

		response.Body = GetResponseBody(score, req.Document)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	response.Body = GetResponseBody(score, req.Document)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderGetScore)
}
