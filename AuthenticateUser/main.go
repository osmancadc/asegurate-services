package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerAuthenticateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
	defer conn.Close()

	//TODO Format Query to use the "?" nomeclature

	results, err := conn.Query(fmt.Sprintf(`SELECT u.document id, CONCAT(p.name," ",p.lastname) name,u.role FROM user u
												INNER JOIN person p on u.document = p.document
												WHERE u.username = '%s' and u.password = '%s'`, data.Username, data.Password))
	if err != nil {
		response.Body = fmt.Sprintf(`{"error_code":"DB01","message":"%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	for results.Next() {
		var document, name, role string

		err = results.Scan(&document, &name, &role)
		if err != nil {
			response.Body = fmt.Sprintf(`{"error_code":"DB02","message":"%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, err
		}

		token, err := GenerateJWT(User{
			UserId: document,
			Name:   name,
			Role:   role,
		})
		if err != nil {
			response.Body = fmt.Sprintf(`{"error_code":"DB03","message":"%s"}`, err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, err
		} else {
			response.Body = fmt.Sprintf(`{"message":"User authenticated","token":"%s"}`, token)
			response.StatusCode = http.StatusOK
			return response, nil
		}
	}

	response.Body = `{"message":"No user found","token":""}`
	response.StatusCode = http.StatusUnauthorized
	return response, nil

}

func main() {
	lambda.Start(HandlerAuthenticateUser)
}
