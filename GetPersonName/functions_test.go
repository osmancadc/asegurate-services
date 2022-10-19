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
				Body:       `{ "message": "success","name":"some_name"}`,
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

func TestGetNameByPhone(t *testing.T) {
	type args struct {
		phone  string
		client lambdaiface.LambdaAPI
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name: `Success Test`,
			args: args{
				phone:  `3123456`,
				client: &MockGetByPhone{},
			},
			wantName: `some_name some_lastname`,
			wantErr:  false,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockGetByPhone{},
			},
			wantName: ``,
			wantErr:  true,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockGetByPhone{},
			},
			wantName: ``,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, err := GetNameByPhone(tt.args.phone, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNameByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("GetNameByPhone() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestGetNameByDocument(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name: `Success Test`,
			args: args{
				document: `123456`,
				client:   &MockGetByDocument{},
			},
			wantName: `some_name some_lastname`,
			wantErr:  false,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockGetByDocument{},
			},
			wantName: ``,
			wantErr:  true,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockGetByDocument{},
			},
			wantName: ``,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, err := GetNameByDocument(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNameByDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("GetNameByDocument() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestGetNameFromProvider(t *testing.T) {
	type args struct {
		documentType string
		document     string
		client       lambdaiface.LambdaAPI
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				document:     `123456`,
				documentType: `CC`,
				client:       &MockGetExternal{},
			},
			want:    true,
			want1:   `some_name some_lastname`,
			wantErr: false,
		},
		{
			name: `Success Test - No Data`,
			args: args{
				client: &MockGetExternal{},
			},
			want:    false,
			want1:   ``,
			wantErr: false,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				documentType: `CC`,
				client:       &MockGetExternal{},
			},
			want:    false,
			want1:   ``,
			wantErr: true,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				documentType: `CC`,
				client:       &MockGetExternal{},
			},
			want:    false,
			want1:   ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetNameFromProvider(tt.args.documentType, tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNameFromProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNameFromProvider() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetNameFromProvider() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetNameFromDatabase(t *testing.T) {
	type args struct {
		dataType  string
		dataValue string
		client    lambdaiface.LambdaAPI
	}
	tests := []struct {
		name      string
		args      args
		wantFound bool
		wantName  string
		wantErr   bool
	}{
		{
			name: `Success Tests - Using CC`,
			args: args{
				dataType:  `CC`,
				dataValue: `123456`,
				client:    &MockGetDatabase{},
			},
			wantFound: true,
			wantName:  `some_name some_lastname`,
			wantErr:   false,
		},
		{
			name: `Success Tests - Using PHONE`,
			args: args{
				dataType:  `PHONE`,
				dataValue: `312345`,
				client:    &MockGetDatabase{},
			},
			wantFound: true,
			wantName:  `some_name some_lastname`,
			wantErr:   false,
		},
		{
			name:      `Success Tests - No Data`,
			args:      args{},
			wantFound: false,
			wantName:  ``,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFound, gotName, err := GetNameFromDatabase(tt.args.dataType, tt.args.dataValue, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNameFromDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetNameFromDatabase() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
			if gotName != tt.wantName {
				t.Errorf("GetNameFromDatabase() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
