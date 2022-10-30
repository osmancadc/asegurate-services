package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var getComments = 1

// GetByPhone Mock
type MockGetComments struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetComments) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getComments {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"comments":[{"author":"some_author","photo":"some_photo","comment":"some_comment","score":0}]}`,
		})
		err = nil
		getComments += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{"message":"some_error"}`,
		})
		err = nil
		getComments += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{})
		err = errors.New(`some_error`)
		getComments += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}
