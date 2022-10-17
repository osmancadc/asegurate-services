package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

var mainTestNumber = 1

// InsertUser Mock
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
			Body:       `{"name":"some_name","email":"some_email","phone":"some_phone","photo":"some_photo","gender":"some_gender"}`,
		})
		err = nil
		mainTestNumber += 1

	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

func TestHandlerGetUserData(t *testing.T) {

	OldGetClient := GetClient
	defer func() { GetClient = OldGetClient }()

	GetClient = func() lambdaiface.LambdaAPI {
		return &MockMain{}
	}

	pathParameter := make(map[string]string)
	pathParameter[`document`] = `123456`

	type args struct {
		req events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: `Success Test HandlerGetUserData`,
			args: args{
				req: events.APIGatewayProxyRequest{
					PathParameters: pathParameter,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerGetUserData(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerGetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("HandlerGetUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}
