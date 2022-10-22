package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	str "strings"
	"time"

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

func SuccessMessage(score Score) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK

	body, _ := json.Marshal(score)
	response.Body = string(body)

	return
}

func GetAssociatedDocument(phone string, client lambdaiface.LambdaAPI) (document string) {
	payload := GetDocumentInvokePayload(phone)
	response := InvokeResponse{}
	person := Person{}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetAssociatedDocument(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	if response.StatusCode != 200 {
		fmt.Printf(`GetAssociatedDocument(2): %s`, response.Body)
		return
	}

	json.Unmarshal([]byte(response.Body), &person)
	document = person.Document

	return
}

// It takes a document, sends it to the `InternalData` Lambda function, and returns the result
func GetInternalScore(document string, client lambdaiface.LambdaAPI) (score InternalScore, isStored bool, err error) {
	payload := GetInternalScoreInvokePayload(document)
	response := InvokeResponse{}
	responseMessage := ResponseMessage{}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetInternalScore(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	if response.StatusCode != 200 {
		err = json.Unmarshal([]byte(response.Body), &responseMessage)
		fmt.Printf(`GetInternalScore(2): %s`, responseMessage.Message)
		return
	}

	json.Unmarshal([]byte(response.Body), &score)
	isStored = true
	return
}

// It takes a document, sends it to the `ExternalData` Lambda function, and returns the result
func GetExternalProccedings(document string, client lambdaiface.LambdaAPI) (proccedings ExternalProccedings, err error) {
	payload := GetExternalProccedingsInvokePayload(document)
	response := InvokeResponse{}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("ExternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetExternalProccedings(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	if response.StatusCode != 200 {
		fmt.Printf(`GetExternalProccedings(2): %s`, response.Body)
		return
	}

	fmt.Println(response.Body)

	json.Unmarshal([]byte(response.Body), &proccedings)

	return
}

func GetStoredScore(document string, client lambdaiface.LambdaAPI) (storedScore Score, daysSinceLastUpdate int) {
	payload := GetStoredReputationInvokePayload(document)
	response := InvokeResponse{}
	storedReputation := StoredReputation{}
	daysSinceLastUpdate = 7

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetStoredScore(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	if response.StatusCode != 200 {
		fmt.Printf(`GetStoredScore(2): %s`, response.Body)
		return
	}

	json.Unmarshal([]byte(response.Body), &storedReputation)

	lastUpdated, err := time.Parse(`2006-01-02 15:04:05`, storedReputation.LastUpdate)
	if err != nil {
		fmt.Printf(`GetStoredScore(1)  %s`, err.Error())
		return
	}

	daysSinceLastUpdate = int(time.Since(lastUpdated).Hours() / 24)
	storedScore = Score{
		Name:       fmt.Sprintf("%s %s", storedReputation.Name, storedReputation.Lastname),
		Reputation: storedReputation.Reputation,
		Gender:     storedReputation.Gender,
		Photo:      storedReputation.Photo,
	}

	return
}

func PredictPersonScore(internalScore InternalScore, externalProccedings ExternalProccedings, client lambdaiface.LambdaAPI) (predictionResponse PredictionResponse, err error) {
	payload := GetPredictScoreInvokePayload(internalScore, externalProccedings)
	response := InvokeResponse{}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("PredictScore"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetExternalProccedings(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	if response.StatusCode != 200 {
		fmt.Printf(`GetExternalProccedings(2): %s`, response.Body)
		return
	}

	json.Unmarshal([]byte(response.Body), &predictionResponse)
	return
}

func CalculateScore(document string, client lambdaiface.LambdaAPI) (isStored bool, score Score, err error) {

	internalScore := InternalScore{}
	externalProccedings := ExternalProccedings{}

	internalScore, isStored, err = GetInternalScore(document, client)
	if err != nil {
		fmt.Printf(`CalculateScore(1): %s`, err.Error())
		return
	}

	score, daysSinceLastUpdate := GetStoredScore(document, client)
	if daysSinceLastUpdate >= 7 || !isStored {
		externalProccedings, err = GetExternalProccedings(document, client)
		if err != nil {
			fmt.Printf(`CalculateScore(2): %s`, err.Error())
			return
		}
	}

	predictionResponse, err := PredictPersonScore(internalScore, externalProccedings, client)

	score.Score = predictionResponse.Score
	score.Document = document
	score.Certified = true
	if daysSinceLastUpdate >= 7 {
		score.Reputation = predictionResponse.Reputation
	}

	return
}

func UpdateSavedReputation(document string, reputation int, client lambdaiface.LambdaAPI) (err error) {
	payload := UpdateSavedReputationInvokePayload(document, reputation)
	response := InvokeResponse{}
	responseMessage := ResponseMessage{}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`UpdateSavedReputation(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	err = json.Unmarshal([]byte(response.Body), &responseMessage)
	if response.StatusCode != 200 {
		fmt.Printf(`UpdateSavedReputation(2): %s`, responseMessage.Message)
		err = errors.New(responseMessage.Message)
		return
	}

	return
}

func GetAssociatedName(document string, client lambdaiface.LambdaAPI) (name, lastname string, err error) {
	response := InvokeResponse{}
	responseMessage := ResponseMessage{}
	person := Person{}

	payload := GetAssociatedNameInvokePayload(document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("ExternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetFromProvider(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	bodyString := str.Replace(string(response.Body), `\`, ``, -1)

	if response.StatusCode != 200 {
		fmt.Printf(`GetFromProvider(2): %s`, response.Body)
		json.Unmarshal([]byte(bodyString), &responseMessage)
		err = errors.New(responseMessage.Message)
		return
	}

	json.Unmarshal([]byte(bodyString), &person)

	name = person.Name
	lastname = person.Lastname

	return
}

func SaveNewReputation(document string, reputation int, client lambdaiface.LambdaAPI) (name, lastname string, err error) {
	response := InvokeResponse{}
	responseMessage := ResponseMessage{}
	name, lastname, err = GetAssociatedName(document, client)
	if err != nil {
		return
	}

	payload := SaveNewReputationInvokePayload(name, lastname, document, reputation)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`SavePerson(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		bodyString := str.Replace(string(response.Body), `\`, ``, -1)
		json.Unmarshal([]byte(bodyString), &responseMessage)

		fmt.Printf(`SavePerson(2): %s`, responseMessage.Message)
		err = errors.New(responseMessage.Message)
		return
	}

	return
}
