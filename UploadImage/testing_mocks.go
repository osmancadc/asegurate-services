package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var updateDatabaseTestNumber = 1

// TestGetStoredScore Mock
type MockUpdateDatabase struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockUpdateDatabase) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch updateDatabaseTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"score":0,"reputation":0}`,
		})
		err = nil
		updateDatabaseTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		updateDatabaseTestNumber += 1
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
