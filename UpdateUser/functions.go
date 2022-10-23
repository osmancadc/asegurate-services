package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	str "strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	invokeLambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func SuccessMessage() (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()
	response.StatusCode = http.StatusOK
	response.Body = `{"message":"image updated successfully"}`
	return
}

func GenerateName(document, date, extension string) string {
	hash := sha512.New()

	hash.Write([]byte(document + ` ` + date))

	name := hash.Sum(nil)

	return fmt.Sprintf("%s.%s", hex.EncodeToString(name), extension)
}

func UploadImage(data []byte, name, document string) (location string, err error) {

	temporalName, _ := SaveTemporalFile(data, name, document)

	file, fileName, _ := GetTemporalFile(temporalName)

	location, err = UploadToS3(file, fileName)

	os.Remove(temporalName)

	return
}

func SaveTemporalFile(data []byte, name, document string) (temporalFileName string, err error) {

	extension := str.Split(name, ".")
	date := time.Now().String()

	name = GenerateName(document, date, extension[1])

	temporalFileName = fmt.Sprintf(`/tmp/%s`, name)
	os.WriteFile(temporalFileName, []byte(data), 0644)

	return
}

func GetTemporalFile(temporalName string) (file io.Reader, fileName string, err error) {
	file, err = os.Open(temporalName)
	if err != nil {
		return
	}

	fileName = filepath.Base(temporalName)

	return
}

func UploadToS3(file io.Reader, fileName string) (location string, err error) {
	bucket := os.Getenv("BUCKET_NAME")
	region := os.Getenv("REGION")
	acl := os.Getenv(`ACL`)

	session, _ := session.NewSession(&aws.Config{Region: &region})
	uploader := s3manager.NewUploader(session)

	uploadParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &fileName,
		Body:   file,
		ACL:    &acl,
	}

	result, err := uploader.Upload(uploadParams)
	if err != nil {
		return
	}

	location = result.Location

	return
}

func UpdateDatabase(route, document, email, phone string, client lambdaiface.LambdaAPI) (err error) {
	payload := UpdateSavedReputationInvokePayload(document, route, email, phone)
	response := InvokeResponse{}
	responseMessage := ResponseMessage{}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`UpdateDatabase(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	err = json.Unmarshal([]byte(response.Body), &responseMessage)

	if response.StatusCode != 200 {
		fmt.Printf(`UpdateDatabase(2): %s`, responseMessage.Message)
		err = errors.New(responseMessage.Message)

		return
	}

	return
}
