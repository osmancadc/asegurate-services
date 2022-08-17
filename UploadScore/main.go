package UploadScore

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

/*
* Little change to check ,the linter and deploy functionality V24
 */
func HanderUploadScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	response.Body = fmt.Sprintf(`{"message":"Mucho gusto %s %s"}`, reqBody.Name, reqBody.LastName)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderUploadScore)
}
