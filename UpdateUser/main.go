package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerUploadScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var reqBody RequestBody

	client := GetClient()

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		return ErrorMessage(err)
	}

	if reqBody.Image != `` {
		decoded, err := base64.StdEncoding.DecodeString(reqBody.Image)
		if err != nil {
			fmt.Println("Errror 1", err.Error())
		}

		route, err := UploadImage(decoded, reqBody.Name, reqBody.Document)
		if err != nil {
			return ErrorMessage(err)
		}

		err = UpdatePhotoDatabase(route, reqBody.Document, client)
		if err != nil {
			return ErrorMessage(err)
		}
	}

	UpdateUserDatabase(reqBody.Document, reqBody.Email, reqBody.Phone, client)

	return SuccessMessage()
}

func main() {
	lambda.Start(HandlerUploadScore)
}
