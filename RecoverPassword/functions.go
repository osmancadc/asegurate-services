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

func SuccessMessage(email string) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()

	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{"email":"%s"}`, email)

	return
}

func GetByDocumentInvokePayload(document string) (payload []byte) {
	getUserBody, _ := json.Marshal(InvokeBody{
		Action: `getAccountdata`,
		GetByDocument: GetByDocumentBody{
			Document: document,
			Fields:   []string{`email`},
		},
	})

	body := InvokePayload{
		Body: string(getUserBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetEmailByDocument(document string, client lambdaiface.LambdaAPI) (email string, err error) {
	response := InvokeResponse{}
	responseBody := ResponseBody{}
	messageBody := MessageBody{}

	payload := GetByDocumentInvokePayload(document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetEmailByDocument(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	bodyString := str.Replace(string(response.Body), `\`, ``, -1)

	if response.StatusCode != 200 {
		fmt.Printf(`GetEmailByDocument(2): %s`, response.Body)
		json.Unmarshal([]byte(bodyString), &messageBody)
		err = errors.New(messageBody.Message)
		return
	}

	json.Unmarshal([]byte(bodyString), &responseBody)
	email = responseBody.Email

	return
}
