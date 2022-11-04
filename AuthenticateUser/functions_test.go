package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	_ "github.com/go-sql-driver/mysql"
)

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

func TestErrorMessage(t *testing.T) {
	type args struct {
		functionError error
		statusCode    int
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
				statusCode:    500,
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
			gotResponse, err := ErrorMessage(tt.args.functionError, tt.args.statusCode)
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

func TestSuccessMessage(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				token: `some_token`,
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"User authenticated","token":"some_token","expiresIn":3600}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := SuccessMessage(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuccessMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("SuccessMessage() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	type args struct {
		document string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				document: `123456`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateJWT(tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetUserPassword(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name         string
		args         args
		wantPassword string
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				client: &MockGetUserPassword{},
			},
			wantPassword: `some_password`,
			wantErr:      false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockGetUserPassword{},
			},
			wantPassword: ``,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPassword, err := GetUserPassword(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPassword != tt.wantPassword {
				t.Errorf("GetUserPassword() = %v, want %v", gotPassword, tt.wantPassword)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	type args struct {
		requestBody RequestBody
		client      lambdaiface.LambdaAPI
	}
	tests := []struct {
		name        string
		args        args
		wantFound   bool
		wantIsValid bool
		wantErr     bool
	}{
		{
			name: `Success Test - User Authenticated`,
			args: args{
				requestBody: RequestBody{
					Password: `some_password`,
				},
				client: &MockValidateUser{},
			},
			wantFound:   true,
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: `Success Test - No User Found`,
			args: args{
				requestBody: RequestBody{},
				client:      &MockValidateUser{},
			},
			wantFound:   false,
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: `Success Test - User Unauthorized`,
			args: args{
				requestBody: RequestBody{
					Password: `some_wrong_password`,
				},
				client: &MockValidateUser{},
			},
			wantFound:   true,
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				requestBody: RequestBody{},
				client:      &MockValidateUser{},
			},
			wantFound:   false,
			wantIsValid: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFound, gotIsValid, err := ValidateUser(tt.args.requestBody, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFound != tt.wantFound {
				t.Errorf("ValidateUser() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("ValidateUser() gotIsValid = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}
