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

func SuccessMessage(comments Comments) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()

	response.StatusCode = http.StatusOK
	commentsJson, _ := json.Marshal(comments)
	response.Body = string(commentsJson)

	return
}

func GetCommentsInvokePayload(document string) (payload []byte) {
	getCommentsBody, _ := json.Marshal(InvokeBody{
		Action: `getComments`,
		GetByDocument: GetByDocumentBody{
			Document: document,
		},
	})

	body := InvokePayload{
		Body: string(getCommentsBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetCommentsFromDatabase(document string, client lambdaiface.LambdaAPI) (comments Comments, err error) {
	response := InvokeResponse{}
	messageBody := MessageBody{}

	payload := GetCommentsInvokePayload(document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetCommentsFromDatabase(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	bodyString := str.Replace(string(response.Body), `\`, ``, -1)

	if response.StatusCode != 200 {
		fmt.Printf(`GetCommentsFromDatabase(2): %s`, response.Body)
		json.Unmarshal([]byte(bodyString), &messageBody)
		err = errors.New(messageBody.Message)
		return
	}

	json.Unmarshal([]byte(bodyString), &comments)

	return
}
