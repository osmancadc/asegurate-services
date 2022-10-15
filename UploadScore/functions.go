package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

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

func GetUploadInvokePayload(data RequestBody) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{Action: `insertScore`, InsertData: data})

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetFindUserInvokePayload(data RequestBody) (payload []byte) {
	uploadBody, _ := json.Marshal(
		InvokeBody{
			Action: `getByPhone`,
			FindByPhoneData: FindByPhoneBody{
				Phone: data.Objective,
			},
		},
	)

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func UploadScore(data RequestBody, client lambdaiface.LambdaAPI) (err error) {

	payload := GetUploadInvokePayload(data)
	if err != nil {
		fmt.Printf(`UploadScore(1): %s`, err.Error())
		return err
	}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`UploadScore(2): %s`, err.Error())
		return err
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		fmt.Printf(`UploadScore(3): %s`, response.Body)
		return errors.New(response.Body)
	}

	return
}

func FindUserByPhone(data RequestBody, client lambdaiface.LambdaAPI) (request RequestBody, err error) {

	payload := GetFindUserInvokePayload(data)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`FindUserByPhone(1): %s`, err.Error())
		return
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		fmt.Printf(`FindUserByPhone(2): %s`, response.Body)
		return
	}

	userByPhone := FindByPhoneResponseBody{}

	err = json.Unmarshal([]byte(response.Body), &userByPhone)

	request = RequestBody{
		Author:    data.Author,
		Type:      `CC`,
		Objective: userByPhone.Document,
		Score:     data.Score,
		Comments:  data.Comments,
	}

	return
}
