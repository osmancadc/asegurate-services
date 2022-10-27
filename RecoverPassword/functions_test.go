package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
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

func TestSuccessMessage(t *testing.T) {
	type args struct {
		name string
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
				name: `some_name`,
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"email":"some_name"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := SuccessMessage(tt.args.name)
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

func TestGetEmailByDocument(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name      string
		args      args
		wantEmail string
		wantErr   bool
	}{
		{
			name: `Success Test`,
			args: args{
				client: &MockGetEmail{},
			},
			wantEmail: `some@email.com`,
			wantErr:   false,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockGetEmail{},
			},
			wantErr: true,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockGetEmail{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := GetEmailByDocument(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEmailByDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEmail != tt.wantEmail {
				t.Errorf("GetEmailByDocument() = %v, want %v", gotEmail, tt.wantEmail)
			}
		})
	}
}
