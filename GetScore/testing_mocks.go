package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var getAssociatedDocumentTestNumber = 1
var getInternalScoreTestNumber = 1
var getExternalProceedingsTestNumber = 1
var getStoredScoreTestNumber = 1
var predictScoreTestNumber = 1
var calculateScoreTestNumber = 1
var updateTestNumber = 1
var getAssociatedNameTestNumber = 1
var createScoreTestNumber = 1

// GetAssociatedDocument Mock
type MockGetAssociatedDocument struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetAssociatedDocument) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getAssociatedDocumentTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"document":"123456"}`,
		})
		err = nil
		getAssociatedDocumentTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getAssociatedDocumentTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getAssociatedDocumentTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestGetInternalScore Mock
type MockGetInternalScore struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetInternalScore) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getInternalScoreTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"score":0,"positive_scores":0,"negative_scores":0,"average_60_days":0}`,
		})
		err = nil
		getInternalScoreTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getInternalScoreTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getInternalScoreTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestGetExternalProceedings Mock
type MockGetExternalProceedings struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetExternalProceedings) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getExternalProceedingsTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"formal_complaints":0,"recent_complain_year":0,"five_years_amount":0}`,
		})
		err = nil
		getExternalProceedingsTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getExternalProceedingsTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getExternalProceedingsTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestGetStoredScore Mock
type MockGetStoredScore struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetStoredScore) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getStoredScoreTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname","reputation":0,"last_update":"2006-01-02 15:04:05"}`,
		})
		err = nil
		getStoredScoreTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getStoredScoreTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getStoredScoreTestNumber += 1
	case 4:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{}`,
		})
		err = nil
		getStoredScoreTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestGetStoredScore Mock
type MockPredictScore struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockPredictScore) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch predictScoreTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"score":0,"reputation":0}`,
		})
		err = nil
		predictScoreTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		predictScoreTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		predictScoreTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestCalculateScore Mock
type MockCalculateScore struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockCalculateScore) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch calculateScoreTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"score":0,"positive_scores":0,"negative_scores":0,"average_60_days":0}`,
		})
		err = nil
		calculateScoreTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname","reputation":0,"last_update":"2006-01-02 15:04:05"}`,
		})
		err = nil
		calculateScoreTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"score":0,"reputation":0}`,
		})
		err = nil
		calculateScoreTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestGetStoredScore Mock
type MockUpdate struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockUpdate) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch updateTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"some_message"}`,
		})
		err = nil
		updateTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		updateTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		updateTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestGetStoredScore Mock
type MockAssociatedName struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockAssociatedName) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getAssociatedNameTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		getAssociatedNameTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getAssociatedNameTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getAssociatedNameTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

// TestGetStoredScore Mock
type MockCreateScore struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockCreateScore) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch createScoreTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		createScoreTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"message":"some_message"}`,
		})
		err = nil
		createScoreTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		createScoreTestNumber += 1
	case 4:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		createScoreTestNumber += 1
	case 5:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","lastname":"some_lastname"}`,
		})
		err = nil
		createScoreTestNumber += 1
	case 6:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		createScoreTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}
