package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerUploadScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	fmt.Printf("Request: %v", reqBody)

	authorId, err := GetAuthorId(conn, reqBody.Author)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if reqBody.Type == `CC` {
		fmt.Printf("Entro a CC")
		err = UploadScoreDocument(conn, authorId, reqBody.Score, reqBody.Value, reqBody.Comments)
		if err != nil {
			response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}
	} else {
		fmt.Printf("Entro a PHONE")
		err = UploadScorePhone(conn, authorId, reqBody.Score, reqBody.Value, reqBody.Comments)
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
	lambda.Start(HandlerUploadScore)
}
