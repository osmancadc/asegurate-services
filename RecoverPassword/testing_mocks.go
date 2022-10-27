package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var getEmailTestNumber = 1
var mainTestNumber = 1

// MockGetEmail
type MockGetEmail struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetEmail) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getEmailTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"email":"some@email.com"}`,
		})
		err = nil
		getEmailTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getEmailTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getEmailTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// MockMain
type MockMain struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockMain) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch mainTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"email":"some@email.com"}`,
		})
		err = nil
		mainTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		mainTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}
