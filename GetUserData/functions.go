package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	str "strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	invokeLambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	_ "github.com/go-sql-driver/mysql"
)

var GetClient = func() lambdaiface.LambdaAPI {
	region := os.Getenv(`REGION`)
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	return invokeLambda.New(sess, &aws.Config{Region: aws.String(region)})
}

func SetResponseHeaders() (response events.APIGatewayProxyResponse) {
	response = events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}
	return
}

func ErrorMessage(functionError error) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()

	response.StatusCode = http.StatusInternalServerError
	response.Body = fmt.Sprintf(`{"message":"%s"}`, functionError.Error())

	return
}

func GetUserDataInvokePayload(document string) (payload []byte) {
	getUserBody, _ := json.Marshal(InvokeBody{
		Action: `getAccountdata`,
		GetUserData: GetByDocumentBody{
			Document: document,
		}})

	body := InvokePayload{
		Body: string(getUserBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetUserData(document string, client lambdaiface.LambdaAPI) (User, error) {
	user := User{
		Document: document,
	}

	payload := GetUserDataInvokePayload(document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetUserData(1): %s`, err.Error())
		return user, err
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		fmt.Printf(`GetUserData(2): %s`, response.Body)
		return user, errors.New(`error obteniendo datos, intentalo de nuevo`)
	}

	bodyString := str.Replace(string(response.Body), `\`, ``, -1)
	json.Unmarshal([]byte(bodyString), &user)

	return user, nil
}
