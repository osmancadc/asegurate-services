package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerGetPersonName(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	client := GetClient()

	dataType := req.PathParameters["type"]
	dataValue := req.PathParameters["value"]

	found, name := GetNameFromDatabase(dataType, dataValue, client)

	if found {
		return SuccessMessage(name)
	}

	found, name, err := GetNameFromProvider(dataType, dataValue, client)
	if err != nil {
		return ErrorMessage(err)
	}

	if found {
		return SuccessMessage(name)
	}

	return ErrorMessage(errors.New(`el usuario no se pudo encontrar`))
}

func main() {
	lambda.Start(HandlerGetPersonName)
}
