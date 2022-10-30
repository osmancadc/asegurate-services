package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerGetPersonName(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	client := GetClient()

	document := req.PathParameters["document"]

	comments, err := GetCommentsFromDatabase(document, client)
	if err != nil {
		return ErrorMessage(err)
	}

	return SuccessMessage(comments)

}

func main() {
	lambda.Start(HandlerGetPersonName)
}
