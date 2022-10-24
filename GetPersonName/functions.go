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

func SuccessMessage(name string) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()

	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{ "message": "success","name":"%s"}`, name)

	return
}

func GetExternalInvokePayload(documentType, document string) (payload []byte) {
	getNameBody, _ := json.Marshal(InvokeBody{
		Action: `getPersonName`,
		GetExternalBody: GetExternalBody{
			Document:     document,
			DocumentType: `CC`,
		},
	})

	body := InvokePayload{
		Body: string(getNameBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetByDocumentInvokePayload(document string) (payload []byte) {
	getUserBody, _ := json.Marshal(InvokeBody{
		Action: `getPersonByDocument`,
		GetByDocument: GetByDocumentBody{
			Document: document,
			Fields:   []string{`name`, `lastname`},
		},
	})

	body := InvokePayload{
		Body: string(getUserBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetByPhoneInvokePayload(phone string) (payload []byte) {
	getUserBody, _ := json.Marshal(InvokeBody{
		Action: `getNameByPhone`,
		GetByPhone: GetByPhoneBody{
			Phone: phone,
		},
	})

	body := InvokePayload{
		Body: string(getUserBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetNameByPhone(phone string, client lambdaiface.LambdaAPI) (name string, err error) {
	response := InvokeResponse{}
	responseBody := ResponseBody{}
	messageBody := MessageBody{}

	payload := GetByPhoneInvokePayload(phone)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetNameByPhone(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	bodyString := str.Replace(string(response.Body), `\`, ``, -1)

	if response.StatusCode != 200 {
		fmt.Printf(`GetNameByPhone(2): %s`, response.Body)
		json.Unmarshal([]byte(bodyString), &messageBody)
		err = errors.New(messageBody.Message)
		return
	}

	err = json.Unmarshal([]byte(response.Body), &responseBody)
	name = fmt.Sprintf(`%s %s`, responseBody.Name, responseBody.Lastname)

	return
}

func GetNameByDocument(document string, client lambdaiface.LambdaAPI) (name string, err error) {
	response := InvokeResponse{}
	responseBody := ResponseBody{}
	messageBody := MessageBody{}

	payload := GetByDocumentInvokePayload(document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetNameByDocument(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	bodyString := str.Replace(string(response.Body), `\`, ``, -1)

	if response.StatusCode != 200 {
		fmt.Printf(`GetNameByDocument(2): %s`, response.Body)
		json.Unmarshal([]byte(bodyString), &messageBody)
		err = errors.New(messageBody.Message)
		return
	}

	json.Unmarshal([]byte(bodyString), &responseBody)
	name = fmt.Sprintf(`%s %s`, responseBody.Name, responseBody.Lastname)

	return
}

func GetNameFromDatabase(dataType, dataValue string, client lambdaiface.LambdaAPI) (bool, string) {
	if dataType == `CC` {
		name, err := GetNameByDocument(dataValue, client)
		if err != nil {
			fmt.Printf(`GetNameFromDatabase(1): %s`, err.Error())
			return false, ``
		}
		return true, name
	} else if dataType == `PHONE` {
		name, err := GetNameByPhone(dataValue, client)
		if err != nil {
			fmt.Printf(`GetNameFromDatabase(2): %s`, err.Error())
			return false, ``
		}
		return true, name
	}

	return false, ""
}

func GetNameFromProvider(documentType, document string, client lambdaiface.LambdaAPI) (bool, string, error) {
	if documentType == `CC` {
		response := InvokeResponse{}
		responseBody := ResponseBody{}
		messageBody := MessageBody{}

		payload := GetExternalInvokePayload(documentType, document)

		result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("ExternalData"), Payload: payload})
		if err != nil {
			fmt.Printf(`GetFromProvider(1): %s`, err.Error())
			return false, ``, err
		}

		json.Unmarshal(result.Payload, &response)
		bodyString := str.Replace(string(response.Body), `\`, ``, -1)

		if response.StatusCode != 200 {
			fmt.Printf(`GetFromProvider(2): %s`, response.Body)
			json.Unmarshal([]byte(bodyString), &messageBody)
			return false, ``, errors.New(messageBody.Message)
		}

		json.Unmarshal([]byte(bodyString), &responseBody)

		name := fmt.Sprintf("%s %s", responseBody.Name, responseBody.Lastname)

		return true, name, nil
	}

	return false, ``, nil
}
