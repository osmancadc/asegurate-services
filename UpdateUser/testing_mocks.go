package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var updatePhotoDatabaseTestNumber = 1
var updateUserDatabaseTestNumber = 1

// UpdatePhotoDatabase Mock
type MockUpdatePhotoDatabase struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockUpdatePhotoDatabase) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch updatePhotoDatabaseTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"score":0,"reputation":0}`,
		})
		err = nil
		updatePhotoDatabaseTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		updatePhotoDatabaseTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// UserDatabase Mock
type MockUpdateUserDatabase struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockUpdateUserDatabase) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch updateUserDatabaseTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"score":0,"reputation":0}`,
		})
		err = nil
		updateUserDatabaseTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		updateUserDatabaseTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}
