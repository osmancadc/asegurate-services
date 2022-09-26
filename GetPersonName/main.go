package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HanderGetPersonName(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	dataType := req.PathParameters["type"]
	dataValue := req.PathParameters["value"]

	fmt.Printf(`%s -> %s`, dataType, dataValue)

	conn := ConnectDatabase()
	defer conn.Close()

	found, name, err := GetFromDatabase(conn, dataType, dataValue)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if found {
		response.Body = fmt.Sprintf(`{ "message": "success","name":"%s"}`, name)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	found, name, err = GetFromProvider(dataType, dataValue)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	if found {
		response.Body = fmt.Sprintf(`{ "message": "success","name":"%s"}`, name)
		response.StatusCode = http.StatusOK
		return response, nil
	}

	response.Body = `{ "message": "El usuario no se pudo encontrar"}`
	response.StatusCode = http.StatusInternalServerError
	return response, nil

}

func main() {
	lambda.Start(HanderGetPersonName)
}
