package main

import (
	"encoding/json"
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
	jwt "github.com/golang-jwt/jwt/v4"

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

func ErrorMessage(functionError error, statusCode int) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()

	response.StatusCode = statusCode
	response.Body = fmt.Sprintf(`{"message":"%s"}`, functionError.Error())

	return
}

func SuccessMessage(token string) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = fmt.Sprintf(`{"message":"User authenticated","token":"%s","expiresIn":3600}`, token)

	return
}

func GetUserPasswordInvokePayload(document string) (payload []byte) {
	getNameBody, _ := json.Marshal(InvokeBody{
		Action: `getUserByDocument`,
		GetByDocument: GetByDocumentBody{
			Document: document,
			Fields:   []string{`password`},
		},
	})

	body := InvokePayload{
		Body: string(getNameBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetUserPassword(document string, client lambdaiface.LambdaAPI) (password string, err error) {
	response := InvokeResponse{}
	responseBody := ResponseBody{}

	payload := GetUserPasswordInvokePayload(document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetUserPassword(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	bodyString := str.Replace(string(response.Body), `\`, ``, -1)

	json.Unmarshal([]byte(bodyString), &responseBody)

	password = responseBody.Password

	return
}

func ValidateUser(requestBody RequestBody, client lambdaiface.LambdaAPI) (found, isValid bool, err error) {

	password, err := GetUserPassword(requestBody.Document, client)
	if err != nil {
		return
	}

	if password == `` {
		return
	}

	found = true
	if password != requestBody.Password {
		return
	}
	isValid = true

	return
}

func GenerateJWT(document string) (token string, err error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Hour * 24),
		},
		ID: document,
	}
	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, _ = tokenData.SignedString([]byte("ASEGUR4TE"))

	return
}
