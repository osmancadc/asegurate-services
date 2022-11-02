package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerUploadScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		return ErrorMessage(err, 500)
	}

	client := GetClient()

	if reqBody.Type == `PHONE` {
		reqBody, err = FindUserByPhone(reqBody, client)
		if err != nil {
			return ErrorMessage(err, 500)
		}

		if reqBody.Author == reqBody.Objective {
			return ErrorMessage(errors.New(`can't score yourself`), 405)
		}

		if reqBody.Objective == `` {
			return ErrorMessage(errors.New(`user not found`), 404)
		}
	}

	err = UploadScore(reqBody, client)
	if err != nil {
		return ErrorMessage(err, 500)
	}

	response := SetResponseHeaders()
	response.Body = `{"message":"Score uploaded successfully"}`
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HandlerUploadScore)
}
