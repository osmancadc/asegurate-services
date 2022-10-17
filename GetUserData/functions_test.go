package main

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	_ "github.com/go-sql-driver/mysql"
)

var getUserDataTestNumber = 1

// InsertUser Mock
type MockGetUserData struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockGetUserData) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {

	var payload []byte
	var err error

	switch getUserDataTestNumber {
	case 1:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
			Body:       `{"name":"some_name","email":"some_email","phone":"some_phone","photo":"some_photo","gender":"some_gender"}`,
		})
		err = nil
		getUserDataTestNumber += 1
	case 2:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 200,
		})
		err = errors.New(`some_error`)
		getUserDataTestNumber += 1
	case 3:
		payload, _ = json.Marshal(InvokeResponse{
			StatusCode: 500,
		})
		err = nil
		getUserDataTestNumber += 1
	}

	return &lambda.InvokeOutput{
		Payload: payload,
	}, err
}

func TestGetUserData(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr bool
	}{
		{
			name: `Success Tests `,
			args: args{
				document: `123456`,
				client:   &MockGetUserData{},
			},
			want: User{
				Document: `123456`,
				Name:     `some_name`,
				Email:    `some_email`,
				Phone:    `some_phone`,
				Photo:    `some_photo`,
				Gender:   `some_gender`,
			},
			wantErr: false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockGetUserData{},
			},
			want:    User{},
			wantErr: true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockGetUserData{},
			},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserData(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserData() = %v, want %v", got, tt.want)
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

func TestSetResponseHeaders(t *testing.T) {
	tests := []struct {
		name         string
		wantResponse events.APIGatewayProxyResponse
	}{
		{
			name: `Success Test`,
			wantResponse: events.APIGatewayProxyResponse{
				Headers: map[string]string{
					"Content-Type":                 "application/json",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "POST",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResponse := SetResponseHeaders(); !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("SetResponseHeaders() = %v, want %v", gotResponse, tt.wantResponse)
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
