package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

// Test counter
var checkExistingPersonTestNumber = 1
var GetPersonDataTestNumber = 1
var savePersonTestNumber = 1
var updateGenderTestNumber = 1
var insertPersonTestNumber = 1
var checkExistingUserTestNumber = 1
var saveUserTestNumber = 1
var insertUserTestNumber = 1

// CheckExistingPerson Mock
type MockCheckExistingPerson struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockCheckExistingPerson) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch checkExistingPersonTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","gender":"some_gender"}`,
		})
		err = nil
		checkExistingPersonTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some server error"}`,
		})
		err = nil
		checkExistingPersonTestNumber += 1

	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{}`,
		})
		err = errors.New(`some error`)
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// GetPersonData  Mock
type MockGetPersonData struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetPersonData) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch GetPersonDataTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name", "lastname":"some_lastname", "gender":"some_gender", "is_alive":true}`,
		})
		err = nil
		GetPersonDataTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_message"}`,
		})
		err = nil
		GetPersonDataTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = errors.New(`some_error`)
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// SavePerson Mock
type MockSavePerson struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockSavePerson) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch savePersonTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"success_message"}`,
		})
		err = nil
		savePersonTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"error_message"}`,
		})
		err = errors.New(`some_error`)
		savePersonTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"error_message"}`,
		})
		err = nil
		savePersonTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// UpdateGender Mock
type MockUpdateGender struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockUpdateGender) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch updateGenderTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"success_message"}`,
		})
		err = nil
		updateGenderTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
		})
		err = errors.New(`some_error`)
		updateGenderTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"error_message"}`,
		})
		err = nil
		updateGenderTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// InsertPerson Mock
type MockInsertPerson struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockInsertPerson) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch insertPersonTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","gender":"some_gender"}`,
		})
		err = nil
		insertPersonTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"not_found","gender":"not_found"}`,
		})
		err = nil
		insertPersonTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// CheckExistingUser Mock
type MockCheckExistingUser struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockCheckExistingUser) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch checkExistingUserTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"user already exists"}`,
		})
		err = nil
		checkExistingUserTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"user does not exists"}`,
		})
		err = nil
		checkExistingUserTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
		})
		err = errors.New(`some_error`)
		checkExistingUserTestNumber += 1
	case 4:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
		})
		err = nil
		checkExistingUserTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// SaveUser Mock
type MockSaveUser struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockSaveUser) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch saveUserTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"success_message"}`,
		})
		err = nil
		saveUserTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"error_message"}`,
		})
		err = errors.New(`some_error`)
		saveUserTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"error_message"}`,
		})
		err = nil
		saveUserTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// InsertUser Mock
type MockInsertUser struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockInsertUser) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch insertUserTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"user does not exists"}`,
		})
		err = nil
		insertUserTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"success_message"}`,
		})
		err = nil
		insertUserTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"user already exists"}`,
		})
		err = nil
		insertUserTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}
