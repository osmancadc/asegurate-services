package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Some change
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

	response.Body = fmt.Sprintf(`{
		"name": "Osman Beltran Murcia",
		"document": "1018500888",
		"stars": 4,
		"reputation": 87,
		"score": 75,
		"certified": true,
		"photo": "data:image/jpeg;base64,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"
	}`)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderGetScore)
}
