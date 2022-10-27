package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

func TestHandlerGetPersonName(t *testing.T) {
	OldGetClient := GetClient
	defer func() { GetClient = OldGetClient }()

	GetClient = func() lambdaiface.LambdaAPI {
		return &MockMain{}
	}

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
			name: `Success Test`,
			args: args{
				events.APIGatewayProxyRequest{
					Body: `{"document":"123456"}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"email":"some@email.com"}`,
			},
			wantErr: false,
		},
		{
			name: `Error Test - Bad Request`,
			args: args{
				events.APIGatewayProxyRequest{
					Body: `{"document":123}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       ``,
			},
			wantErr: false,
		},
		{
			name: `Error Test - No User Found`,
			args: args{
				events.APIGatewayProxyRequest{
					Body: `{"document":"123"}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"el usuario no se pudo encontrar"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerGetPersonName(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerGetPersonName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("HandlerGetPersonName() = %v, want %v", got, tt.want)
			}
		})
	}
}
