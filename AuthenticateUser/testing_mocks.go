package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var getUserPasswordTestNumber = 1
var validateUserTestNumber = 1

type MockGetUserPassword struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetUserPassword) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getUserPasswordTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"password":"some_password"}`,
		})
		err = nil
		getUserPasswordTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getUserPasswordTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// GetByDocument Mock
type MockValidateUser struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockValidateUser) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch validateUserTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"password":"some_password"}`,
		})
		err = nil
		validateUserTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{}`,
		})
		err = nil
		validateUserTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"password":"some_password"}`,
		})
		err = nil
		validateUserTestNumber += 1
	case 4:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// // GetExternal Mock
// type MockGetExternal struct {
// 	lambdaiface.LambdaAPI
// }

// func (mlc *MockGetExternal) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

// 	var payload []byte
// 	var err error

// 	switch getExternalTestNumber {
// 	case 1:
// 		payload, _ = json.Marshal(InvokeResponse{
// 			StatusCode: 200,
// 			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
// 		})
// 		err = nil
// 		getExternalTestNumber += 1
// 	case 2:
// 		payload, _ = json.Marshal(InvokeResponse{
// 			StatusCode: 500,
// 			Body:       `{"message":"some_error"}`,
// 		})
// 		err = nil
// 		getExternalTestNumber += 1
// 	case 3:
// 		payload, _ = json.Marshal(InvokeResponse{})
// 		err = errors.New(`some_error`)
// 		getExternalTestNumber += 1
// 	}

// 	return &lambda.InvokeOutput{
// 		Payload: payload,
// 	}, err
// }

// // GetExternal Mock
// type MockGetDatabase struct {
// 	lambdaiface.LambdaAPI
// }

// func (mlc *MockGetDatabase) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

// 	var payload []byte
// 	var err error

// 	switch getDatabaseTestNumber {
// 	case 1:
// 		payload, _ = json.Marshal(InvokeResponse{
// 			StatusCode: 200,
// 			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
// 		})
// 		err = nil
// 		getDatabaseTestNumber += 1
// 	case 2:
// 		payload, _ = json.Marshal(InvokeResponse{
// 			StatusCode: 200,
// 			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
// 		})
// 		err = nil
// 		getDatabaseTestNumber += 1

// 	}

// 	return &lambda.InvokeOutput{
// 		Payload: payload,
// 	}, err
// }

// // Main Mock
// type MockGetPersonNameMain struct {
// 	lambdaiface.LambdaAPI
// }

// func (mlc *MockGetPersonNameMain) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

// 	var payload []byte
// 	var err error

// 	switch mainTestNumber {
// 	case 1:
// 		payload, _ = json.Marshal(InvokeResponse{
// 			StatusCode: 200,
// 			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
// 		})
// 		err = nil
// 		mainTestNumber += 1
// 	}

// 	return &lambda.InvokeOutput{
// 		Payload: payload,
// 	}, err
// }
