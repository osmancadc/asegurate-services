package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var getByPhoneTestNumber = 1
var getByDocumentTestNumber = 1
var getExternalTestNumber = 1
var getDatabaseTestNumber = 1
var mainTestNumber = 1

// GetByPhone Mock
type MockGetByPhone struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetByPhone) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getByPhoneTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		getByPhoneTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getByPhoneTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getByPhoneTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// GetByDocument Mock
type MockGetByDocument struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetByDocument) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getByDocumentTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		getByDocumentTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getByDocumentTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getByDocumentTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// GetExternal Mock
type MockGetExternal struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetExternal) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getExternalTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		getExternalTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getExternalTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getExternalTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// GetExternal Mock
type MockGetDatabase struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetDatabase) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getDatabaseTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		getDatabaseTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		getDatabaseTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// Main Mock
type MockGetPersonNameMain struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetPersonNameMain) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch mainTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		mainTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}
