package main

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	_ "github.com/go-sql-driver/mysql"
)

type MockFindByIdClient struct {
	lambdaiface.LambdaAPI
}

var testNumber = 1

func (mlc *MockFindByIdClient) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	var payload []byte
	var err error

	switch testNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"document":"654321"}`,
		})
		err = nil
		testNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{}`,
		})
		err = errors.New(`some error`)
		testNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
			Body:       `{}`,
		})
		err = nil
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err

}

type MockUploadScore struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockUploadScore) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	payload, _ := json.Marshal(InvokeResponse{
		StatusCode: 500,
		Body:       `{"message":"some error"}`,
	})

	return &lambda.InvokeOutput{
		Payload: payload,
	}, nil
}

func TestFindUserByPhone(t *testing.T) {
	os.Setenv(`REGION`, `some-region-1`)
	mockLambda := &MockFindByIdClient{}
	type args struct {
		data   RequestBody
		client lambdaiface.LambdaAPI
	}
	tests := []struct {
		name        string
		args        args
		wantRequest RequestBody
		wantErr     bool
	}{
		{
			name: "Success Test",
			args: args{
				client: mockLambda,
				data: RequestBody{
					Author:    `123`,
					Type:      `some_type`,
					Objective: `123456`,
					Score:     100,
					Comments:  `test comment`,
				},
			},
			wantRequest: RequestBody{
				Author:    `123`,
				Type:      `CC`,
				Objective: `654321`,
				Score:     100,
				Comments:  `test comment`,
			},
			wantErr: false,
		},
		{
			name: "Error test - Invoke Error",
			args: args{
				client: mockLambda,
				data:   RequestBody{},
			},
			wantRequest: RequestBody{},
			wantErr:     true,
		},
		{
			name: "Success test - Response Error",
			args: args{
				client: mockLambda,
				data:   RequestBody{},
			},
			wantRequest: RequestBody{},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRequest, err := FindUserByPhone(tt.args.data, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRequest, tt.wantRequest) {
				t.Errorf("FindUserByPhone() = %v, want %v", gotRequest, tt.wantRequest)
			}
		})
	}
}

func TestUploadScore(t *testing.T) {
	type args struct {
		data   RequestBody
		client lambdaiface.LambdaAPI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error test bad response from the lambda",
			args: args{
				data: RequestBody{
					Author:    `123456`,
					Type:      `some type`,
					Objective: `654321`,
					Score:     50,
					Comments:  `some comment`,
				},
				client: &MockUploadScore{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UploadScore(tt.args.data, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("UploadScore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	type args struct {
		functionError error
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test",
			args: args{
				functionError: errors.New(`some error`),
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some error"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := ErrorMessage(tt.args.functionError)
			if (err != nil) != tt.wantErr {
				t.Errorf("ErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("ErrorMessage() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetClient(t *testing.T) {
	tests := []struct {
		name string
		want lambdaiface.LambdaAPI
	}{
		{
			name: `Single test`,
			want: &lambda.Lambda{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetClient()
		})
	}
}
