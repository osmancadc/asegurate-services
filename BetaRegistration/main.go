package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//2
func HanderUploadScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var data RequestBody

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	conn := ConnectDatabase()

	query, err := conn.Prepare(fmt.Sprintf(`INSERT INTO dev_asegurate.beta_user 
											(name, age, cellphone, email, smartphone, operative_system, commentary, associate)
											VALUES('%s', %d, '%s', '%s', %t, '%s', '%s', '%s')`,
		data.Name, data.Age, data.Cellphone, data.Email, data.Smartphone, data.OperativeSystem, data.Commentary, data.Associate))
	if err != nil {
		response.Body = fmt.Sprintf(`{"error_code":"DB01","message":"%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	query.Exec()

	response.Body = fmt.Sprintf(`{"message":"Beta user %s registered"}`, data.Name)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HanderUploadScore)
}
