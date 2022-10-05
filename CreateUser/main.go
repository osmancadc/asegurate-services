package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerCreateUser(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	conn, err := ConnectDatabase()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}
	defer conn.Close()

	name, err := InsertPerson(conn, reqBody.Document, reqBody.ExpeditionDate)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	err = InsertUser(conn, reqBody.Email, reqBody.Phone, reqBody.Password, reqBody.Document, reqBody.Role)
	if err != nil {
		response.Body = fmt.Sprintf(`{ "message": "%s"}`, err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.Body = fmt.Sprintf(`{ "message": "user created successfully","name":"%s"}`, name)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(HandlerCreateUser)
}

// package main

// import (
// 	"fmt"
// 	"os"
// )

// func main() {

// 	os.Setenv("DB_USER", "administrator")
// 	os.Setenv("DB_PASSWORD", "35Yw!8uO5v5g")
// 	os.Setenv("DB_HOST", "dev-asegurate.cluster-cnaioe8hvyno.us-east-1.rds.amazonaws.com")
// 	os.Setenv("DB_NAME", "dev_asegurate")
// 	conn, err := ConnectDatabase()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	defer conn.Close()
// 	name, exists, err := CheckExistingPerson(conn, `1022372535`)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	fmt.Printf("%s -> %t \n", name, exists)
// }
